#include <X11/Xatom.h>
#include <X11/Xlib.h>
#include <X11/extensions/Xfixes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

Atom clipAtom;
Atom utf8Atom;
Atom compoundAtom;
Atom xselAtom;
Atom incrAtom;
Atom targetsAtom;
Atom multipleAtom;
Atom deleteAtom;
Atom timestampAtom;
Atom nullAtom;
int supportedTargetsLen;
Atom supportedTargets[5];

void sleep(long milliseconds) {
	struct timespec ts;
	ts.tv_sec = milliseconds / 1000;
	ts.tv_nsec = (milliseconds % 1000) * 1000000;

	nanosleep(&ts, &ts);
}

Display * displayInit() {
	Display *display = XOpenDisplay(NULL);
	if (!display) {
		fprintf(stderr, "Can't open X display\n");
		exit(974);
	}

	clipAtom = XInternAtom(display, "CLIPBOARD", False);
	utf8Atom = XInternAtom(display, "UTF8_STRING", True);
	compoundAtom = XInternAtom(display, "COMPOUND_TEXT", False);
	xselAtom = XInternAtom(display, "XSEL_DATA", False);
	incrAtom = XInternAtom(display, "INCR", False);
	targetsAtom = XInternAtom(display, "TARGETS", False);
	multipleAtom = XInternAtom(display, "MULTIPLE", False);
	deleteAtom = XInternAtom(display, "DELETE", False);
	timestampAtom = XInternAtom(display, "TIMESTAMP", False);
	nullAtom = XInternAtom(display, "NULL", False);

	int i = 0;
	supportedTargets[i++] = timestampAtom;
	supportedTargets[i++] = targetsAtom;
	supportedTargets[i++] = utf8Atom;
	//supportedTargets[i++] = incrAtom;
	supportedTargets[i++] = deleteAtom;
	supportedTargetsLen = i;

	return display;
}

Window windowInit(Display *display) {
	int blackPixel = BlackPixel(display, DefaultScreen(display));
	Window rootWindow = XDefaultRootWindow(display);
	Window window = XCreateSimpleWindow(display, rootWindow,
		0, 0, 1, 1, 0, blackPixel, blackPixel);

	return window;
}

void displayClose(Display *display) {
	XCloseDisplay(display);
}

void windowClose(Display *display, Window window) {
	XFlush(display);
	XSync(display, False);
	XDestroyWindow(display, window);
}

Time timestampGet(Display *display, Window window) {
	XSelectInput(display, window, PropertyChangeMask);

	XChangeProperty(display, window, XA_WM_NAME, XA_STRING, 8,
		PropModeAppend, NULL, 0);

	XEvent event;
	while (True) {
		XNextEvent(display, &event);

		if (event.type == PropertyNotify) {
			return event.xproperty.time;
		} else {
			fprintf(stderr, "Unexpected timestamp event %d\n", event.type);
		}
	}
}

void clipboardWait() {
	Display *display = displayInit();

	XFlush(display);
	XSync(display, True);

	Window window = windowInit(display);
	int eventBase;
	int errorBase;

	if (!XFixesQueryExtension(display, &eventBase, &errorBase)) {
		fprintf(stderr, "XFixes extension not found\n");
		exit(971);
	}

	XFixesSelectSelectionInput(display, window,
		clipAtom, XFixesSetSelectionOwnerNotifyMask);

	XEvent event;
	while (True) {
		XNextEvent(display, &event);

		if (event.type == eventBase + XFixesSelectionNotify &&
			((XFixesSelectionNotifyEvent*)&event)->selection == clipAtom) {

			break;
		}
	}

	XDeleteProperty(
		event.xselection.display,
		event.xselection.requestor,
		event.xselection.property);

	windowClose(display, window);
	displayClose(display);
}

static Bool appendBuffer(XSelectionEvent *selectEvent, unsigned char **buffer,
		unsigned long *offset, unsigned long *alloc) {
	Atom returnAtom;
	int returnFormat;
	unsigned long length;
	unsigned long bytesAfter;
	unsigned char *value;

	XGetWindowProperty(
		selectEvent->display,
		selectEvent->requestor,
		selectEvent->property,
		0L, 1000000, True, AnyPropertyType,
		&returnAtom, &returnFormat, &length, &bytesAfter, &value);

	if (returnAtom != XA_STRING && returnAtom != utf8Atom &&
		returnAtom != compoundAtom) {

		free(*buffer);
		*buffer = NULL;
		return False;
	}

	if (length == 0) {
		if (value != NULL) {
			XFree(value);
		}
		return False;
	}

	if (returnFormat != 8) {
		if (value != NULL) {
			XFree(value);
		}
		return True;
	}

	if (*offset + length + 1 > *alloc) {
		*alloc = *offset + length + 1;
		if ((*buffer = realloc(*buffer, *alloc)) == NULL) {
			exit(972);
		}
	}

	unsigned char *ptr = *buffer + *offset;

	memcpy(ptr, value, length);
	ptr[length] = '\0';
	*offset += length;

	if (value != NULL) {
		XFree(value);
	}

	return True;
}

char * selectionGet(Display *display, Window window, Time timestamp,
	Atom selectionSource, Atom targetAtom) {

	unsigned char *buffer;
	unsigned long bufferSize;

	xselAtom = XInternAtom(display, "XSEL_DATA", False);
	XConvertSelection(display, selectionSource, targetAtom,
		xselAtom, window, timestamp);

	XEvent convertEvent;
	while (True) {
		XNextEvent(display, &convertEvent);

		if (convertEvent.type == SelectionNotify &&
			convertEvent.xselection.selection == selectionSource) {

			break;
		}  else {
//			fprintf(stderr, "Unexpected selection event %d\n",
//				convertEvent.type);
		}
	}

	if (convertEvent.xselection.property == None) {
//		fprintf(stderr, "Get selection conversion refused\n");
		return NULL;
	}

	Atom returnAtom;
	int returnFormat;
	unsigned long length;
	unsigned long bytesAfter;
	unsigned char *value;

	XGetWindowProperty(
		convertEvent.xselection.display,
		convertEvent.xselection.requestor,
		convertEvent.xselection.property,
		0L, 1000000, False, AnyPropertyType,
		&returnAtom, &returnFormat, &length, &bytesAfter, &value);

	if (returnAtom == incrAtom) {
		bufferSize = *(long *)value;

		XSelectInput(
			convertEvent.xselection.display,
			convertEvent.xselection.requestor,
			PropertyChangeMask);

		XDeleteProperty(
			convertEvent.xselection.display,
			convertEvent.xselection.requestor,
			convertEvent.xselection.property);

		buffer = malloc(bufferSize);
		if (buffer == NULL) {
			fprintf(stderr, "Malloc failed\n");
			exit(973);
		}

		unsigned long incrOffset = 0;
		XEvent propEvent;
		while (True) {
			XNextEvent(convertEvent.xselection.display, &propEvent);

			if (propEvent.type == PropertyNotify &&
				propEvent.xproperty.state == PropertyNewValue) {

				Bool more = appendBuffer(&convertEvent.xselection, &buffer,
					&incrOffset, &bufferSize);
				if (!more) {
					break;
				}
			}
		}
	} else if (returnAtom != XA_STRING && returnAtom != utf8Atom &&
		returnAtom != compoundAtom) {

		return NULL;
	} else {
		buffer = (unsigned char *)strdup((const char *)value);
	}

	if (value != NULL) {
		XFree(value);
	}
	XDeleteProperty(
		convertEvent.xselection.display,
		convertEvent.xselection.requestor,
		convertEvent.xselection.property);

	return (char *)buffer;
}

Bool selectionHandle(XEvent event, Time timestamp, Atom targetAtom,
	unsigned char *buffer) {

	XSelectionEvent selEvent;
	XSelectionRequestEvent *selReq = &event.xselectionrequest;
	Bool more = True;

	selEvent.type = SelectionNotify;
	selEvent.display = selReq->display;
	selEvent.requestor = selReq->requestor;
	selEvent.selection = selReq->selection;
	selEvent.target = selReq->target;
	selEvent.time = selReq->time;

	if (selReq->property == None && selReq->target != multipleAtom) {
		selReq->property = selReq->target;
	}

	if (selEvent.time != CurrentTime && selEvent.time < timestamp) {
		selEvent.property = None;
	} else if (selEvent.target == timestampAtom) {
		selEvent.property = selReq->property;

		XChangeProperty(selEvent.display, selEvent.requestor,
			selEvent.property, XA_INTEGER, 32, PropModeReplace,
			(unsigned char *)timestamp, 1);

		return True;
	} else if (selEvent.target == targetsAtom) {
		selEvent.property = selReq->property;

		Atom *targets = malloc(sizeof(supportedTargets));
		memcpy(targets, supportedTargets, sizeof(supportedTargets));

		XChangeProperty(selEvent.display, selEvent.requestor,
			selEvent.property, XA_ATOM, 32, PropModeReplace,
			(unsigned char *)targets, supportedTargetsLen);

		free(targets);
	} else if (selEvent.target == utf8Atom) {
		selEvent.property = selReq->property;

		XChangeProperty(selEvent.display, selEvent.requestor,
			selEvent.property, targetAtom, 8, PropModeReplace,
			buffer, (int)strlen((char *)buffer));

		more = False;
	} else if (selEvent.target == deleteAtom) {
		selEvent.property = selReq->property;

		XChangeProperty(selEvent.display, selEvent.requestor,
			selEvent.property, nullAtom, 0, PropModeReplace,
			NULL, 0);

		more = False;
	} else {
		selEvent.property = None;
	}

	XSendEvent(selEvent.display, selEvent.requestor, False,
		NoEventMask, (XEvent *)&selEvent);
	return more;
}

void selectionSet(Display *display, Window window, Time timestamp,
	Atom selectionSource, Atom targetAtom, unsigned char *buffer) {

	int eventBase;
	int errorBase;
	if (!XFixesQueryExtension(display, &eventBase, &errorBase)) {
		fprintf(stderr, "XFixes extension not found\n");
		exit(971);
	}
	XFixesSelectSelectionInput(display, window,
		clipAtom, XFixesSetSelectionOwnerNotifyMask);

	XSetSelectionOwner(display, selectionSource, window, timestamp);

	Window owner = XGetSelectionOwner(display, selectionSource);
	if (owner != window) {
		fprintf(stderr, "Failed to get ownership\n");
		return;
	}

	utf8Atom = XInternAtom(display, "UTF8_STRING", True);

	Bool ownerEvent = False;
	Bool selectionEvent = False;
	XEvent event;
	while (True) {
		XNextEvent(display, &event);

		if (event.type == SelectionClear &&
			event.xselectionclear.selection == selectionSource) {

			fprintf(stderr, "Lost ownership\n");
			return;
		} else if (event.type == SelectionRequest &&
			event.xselectionrequest.selection == selectionSource) {

			Bool more = selectionHandle(event, timestamp, targetAtom, buffer);
			if (!more) {
				break;
			}
		}
	}

	timestamp = timestampGet(display, window);
}

char * clipboardGet() {
	Display *display = displayInit();
	Window window = windowInit(display);
	Time timestamp = timestampGet(display, window);

	char *buffer = selectionGet(display, window, timestamp,
		clipAtom, utf8Atom);
	if (buffer == NULL) {
		buffer = selectionGet(display, window, timestamp,
			clipAtom, XA_STRING);
	}

	windowClose(display, window);
	displayClose(display);

	return buffer;
}

void clipboardSet(char *data) {
	Display *display = displayInit();
	Window window = windowInit(display);
	Time timestamp = timestampGet(display, window);

	selectionSet(display, window, timestamp,
		clipAtom, utf8Atom, (unsigned char *)data);

	windowClose(display, window);
	displayClose(display);
}
