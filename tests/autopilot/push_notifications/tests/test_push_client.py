# -*- Mode: Python; coding: utf-8; indent-tabs-mode: nil; tab-width: 4 -*-
# Copyright 2014 Canonical
#
# This program is free software: you can redistribute it and/or modify it
# under the terms of the GNU General Public License version 3, as published
# by the Free Software Foundation.

"""Tests for Push Notifications client"""

import time
from push_notifications.tests import PushNotificationTestBase


class TestPushClient(PushNotificationTestBase):
    """
    Test cases for push notifications
    """

    def test_broadcast_push_notification_screen_off(self):
        """
        Send a push message whilst the device's screen is turned off
        Notification should still be displayed when it is turned on
        """
        # Turn display off
        self.press_power_button()
        # send message
        self.send_valid_push_message()
        # wait before turning screen on
        time.sleep(2)
        # Turn display on
        self.press_power_button()
        # wait to make sure dialog is displayed
        time.sleep(2)
        self.validate_and_dismiss_notification_dialog(
            self.DEFAULT_DISPLAY_MESSAGE)

    def test_broadcast_push_notification_locked_greeter(self):
        """
        Positive test case to send a valid broadcast push notification
        to the client and validate that a notification message is displayed
        whist the greeter screen is displayed
        """
        # Assumes greeter starts in locked state
        self.send_valid_push_message()
        self.validate_and_dismiss_notification_dialog(
            self.DEFAULT_DISPLAY_MESSAGE)

    def test_broadcast_push_notification(self):
        """
        Positive test case to send a valid broadcast push notification
        to the client and validate that a notification message is displayed
        """
        # Assumes greeter starts in locked state
        self.unlock_greeter()
        # send message
        self.send_valid_push_message()
        self.validate_and_dismiss_notification_dialog(
            self.DEFAULT_DISPLAY_MESSAGE)

    def test_expired_broadcast_push_notification(self):
        """
        Send an expired broadcast notification message to server
        """
        self.unlock_greeter()
        msg_data = self.create_notification_data_copy()
        msg_data.inc_build_number()
        msg = self.push_helper.create_push_message(
            data=msg_data.json(),
            expire_after=self.push_helper.get_past_iso_time())
        response = self.push_helper.send_push_broadcast_notification(
            msg.json(),
            self.test_config.server_listener_addr)
        # 400 status is received for an expired message
        self.validate_response(response, expected_status_code=400)
        # validate no notification is displayed
        self.validate_notification_not_displayed()

    def test_older_version_broadcast_push_notification(self):
        """
        Send an old version broadcast notification message to server
        """
        self.unlock_greeter()
        msg_data = self.create_notification_data_copy()
        msg_data.dec_build_number()
        msg = self.push_helper.create_push_message(data=msg_data.json())
        response = self.push_helper.send_push_broadcast_notification(
            msg.json(),
            self.test_config.server_listener_addr)
        self.validate_response(response)
        # validate no notification is displayed
        self.validate_notification_not_displayed()

    def test_equal_version_broadcast_push_notification(self):
        """
        Send an equal version broadcast notification message to server
        """
        self.unlock_greeter()
        msg_data = self.create_notification_data_copy()
        msg = self.push_helper.create_push_message(data=msg_data.json())
        response = self.push_helper.send_push_broadcast_notification(
            msg.json(),
            self.test_config.server_listener_addr)
        self.validate_response(response)
        # validate no notification is displayed
        self.validate_notification_not_displayed()
