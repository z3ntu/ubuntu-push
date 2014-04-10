# -*- Mode: Python; coding: utf-8; indent-tabs-mode: nil; tab-width: 4 -*-
# Copyright 2014 Canonical
#
# This program is free software: you can redistribute it and/or modify it
# under the terms of the GNU General Public License version 3, as published
# by the Free Software Foundation.

"""Tests for Push Notifications client"""

from __future__ import absolute_import

from testtools.matchers import Equals
from autopilot.matchers import Eventually
from autopilot.introspection import dbus

from autopilot.testcase import AutopilotTestCase

from push_notifications.tests import PushNotificationTestBase
from push_notifications.tests import PushNotificationMessage


class TestPushClient(PushNotificationTestBase):
    """ Tests a Push notification can be sent and received """

    def test_get_config(self):
        msg = PushNotificationMessage(expire_after="2015-12-19T16:39:57-08:00")      
        server_add = self.get_push_server_listener_address()
        response = self.send_push_broadcast_notification(server_add, msg.json())
        status = response[0]['status']
        self.assertThat(status, Equals('200'))

