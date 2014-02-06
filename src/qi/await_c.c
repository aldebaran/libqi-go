/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
*/

#include "_cgo_export.h"

void go_future_add_callback(qi_future_t *fut, void *userdata) {
	qi_future_add_callback(fut, (qi_future_callback_t)&go_async_waiter_callback, userdata);
}
