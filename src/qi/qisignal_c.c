/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
*/

#include "_cgo_export.h"

qi_future_t* go_object_signal_connect(qi_object_t *obj, char *signame, void *userdata) {
	return qi_object_signal_connect(obj, signame, (qi_object_signal_callback_t)&go_object_signal_callback, userdata);
}
