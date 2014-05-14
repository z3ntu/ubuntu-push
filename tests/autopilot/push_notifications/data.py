# -*- Mode: Python; coding: utf-8; indent-tabs-mode: nil; tab-width: 4 -*-
# Copyright 2014 Canonical
#
# This program is free software: you can redistribute it and/or modify it
# under the terms of the GNU General Public License version 3, as published
# by the Free Software Foundation.

"""Push-Notifications autopilot data structure classes"""


class PushNotificationMessage:
    """
    Class to hold all the details required for a
    push notification message
    """
    channel = None
    expire_after = None
    data = None

    def __init__(self, channel='system', data='', expire_after=''):
        self.channel = channel
        self.data = data
        self.expire_after = expire_after

    def json(self):
        """
        Return json representation of message
        """
        json_str = '{{"channel":"{0}", "data":{{{1}}}, "expire_on":"{2}"}}'
        return json_str.format(self.channel, self.data, self.expire_after)


class NotificationData:
    """
    Class to represent notification data including
    Device software channel
    Device build number
    Device model
    Device last update
    Data for the notification
    """

    def __init__(self, channel=None, device=None, build_number=None,
        last_update=None, version=None, data=None):
        self.channel = channel
        self.build_number = build_number
        self.device = device
        self.last_update = last_update
        self.version = version
        self.data = data

    def inc_build_number(self):
        """
        Increment build number
        """
        self.build_number = str(int(self.build_number) + 1)

    def dec_build_number(self):
        """
        Decrement build number
        """
        self.build_number = str(int(self.build_number) - 1)

    def json(self):
        """
        Return json representation of info based:
        "IMAGE-CHANNEL/DEVICE-MODEL": [BUILD-NUMBER, CHANNEL-ALIAS]"
        """
        json_str = '"{0}/{1}": [{2}, "{3}"]'
        return json_str.format(self.channel, self.device, self.build_number,
                               self.data)
