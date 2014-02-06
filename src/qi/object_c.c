/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
*/

#include "_cgo_export.h"

unsigned long long go_object_builder_advertise_method(qi_object_builder_t *ob,
													  char *signame,
													  void *callback) {
  //USED(callback);
	return qi_object_builder_advertise_method(ob, signame, (qi_object_method_t)&go_object_call_callback, callback);
}
