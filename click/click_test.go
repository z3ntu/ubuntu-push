/*
 Copyright 2013-2014 Canonical Ltd.

 This program is free software: you can redistribute it and/or modify it
 under the terms of the GNU General Public License version 3, as published
 by the Free Software Foundation.

 This program is distributed in the hope that it will be useful, but
 WITHOUT ANY WARRANTY; without even the implied warranties of
 MERCHANTABILITY, SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR
 PURPOSE.  See the GNU General Public License for more details.

 You should have received a copy of the GNU General Public License along
 with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package click

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "launchpad.net/gocheck"
)

func TestClick(t *testing.T) { TestingT(t) }

type clickSuite struct{}

var _ = Suite(&clickSuite{})

func (cs *clickSuite) TestParseAppId(c *C) {
	app, err := ParseAppId("com.ubuntu.clock_clock")
	c.Assert(err, IsNil)
	c.Check(app.Package, Equals, "com.ubuntu.clock")
	c.Check(app.InPackage("com.ubuntu.clock"), Equals, true)
	c.Check(app.Application, Equals, "clock")
	c.Check(app.Version, Equals, "")
	c.Check(app.Click, Equals, true)
	c.Check(app.Original(), Equals, "com.ubuntu.clock_clock")
	c.Check(fmt.Sprintf("%s", app), Equals, "com.ubuntu.clock_clock")
	c.Check(app.DispatchPackage(), Equals, "com.ubuntu.clock")

	app, err = ParseAppId("com.ubuntu.clock_clock_10")
	c.Assert(err, IsNil)
	c.Check(app.Package, Equals, "com.ubuntu.clock")
	c.Check(app.InPackage("com.ubuntu.clock"), Equals, true)
	c.Check(app.Application, Equals, "clock")
	c.Check(app.Version, Equals, "10")
	c.Check(app.Click, Equals, true)
	c.Check(app.Original(), Equals, "com.ubuntu.clock_clock_10")
	c.Check(fmt.Sprintf("%s", app), Equals, "com.ubuntu.clock_clock_10")
	c.Check(app.Versioned(), Equals, "com.ubuntu.clock_clock_10")
	c.Check(app.Base(), Equals, "com.ubuntu.clock_clock")
	c.Check(app.DesktopId(), Equals, "com.ubuntu.clock_clock_10.desktop")
	c.Check(app.DispatchPackage(), Equals, "com.ubuntu.clock")

	for _, s := range []string{"com.ubuntu.clock_clock_10_4", "com.ubuntu.clock", ""} {
		app, err = ParseAppId(s)
		c.Check(app, IsNil)
		c.Check(err, Equals, ErrInvalidAppId)
	}
}

func (cs *clickSuite) TestVersionedPanic(c *C) {
	app, err := ParseAppId("com.ubuntu.clock_clock")
	c.Assert(err, IsNil)
	c.Check(func() { app.Versioned() }, PanicMatches, `Versioned\(\) on AppId without version/not verified:.*`)
}

func (cs *clickSuite) TestParseAppIdLegacy(c *C) {
	app, err := ParseAppId("_python3.4")
	c.Assert(err, IsNil)
	c.Check(app.Package, Equals, "")
	c.Check(app.InPackage(""), Equals, true)
	c.Check(app.Application, Equals, "python3.4")
	c.Check(app.Version, Equals, "")
	c.Check(app.Click, Equals, false)
	c.Check(app.Original(), Equals, "_python3.4")
	c.Check(app.Versioned(), Equals, "python3.4")
	c.Check(app.Base(), Equals, "python3.4")
	c.Check(app.DesktopId(), Equals, "python3.4.desktop")
	c.Check(app.DispatchPackage(), Equals, "python3.4")

	for _, s := range []string{"_.foo", "_foo/", "_/foo"} {
		app, err = ParseAppId(s)
		c.Check(app, IsNil)
		c.Check(err, Equals, ErrInvalidAppId)
	}
}

func (cs *clickSuite) TestJSON(c *C) {
	for _, appId := range []string{"com.ubuntu.clock_clock", "com.ubuntu.clock_clock_10", "_python3.4"} {
		app, err := ParseAppId(appId)
		c.Assert(err, IsNil, Commentf(appId))
		b, err := json.Marshal(app)
		c.Assert(err, IsNil, Commentf(appId))
		var vapp *AppId
		err = json.Unmarshal(b, &vapp)
		c.Assert(err, IsNil, Commentf(appId))
		c.Check(vapp, DeepEquals, app)
	}
}

func (cs *clickSuite) TestIcon(c *C) {
	app, err := ParseAppId("_python3.4")
	c.Assert(err, IsNil)
	c.Check(app.Icon(), Equals, "/usr/share/pixmaps/python3.4.xpm")
}

func (s *clickSuite) TestUser(c *C) {
	u, err := User()
	c.Assert(err, IsNil)
	c.Assert(u, NotNil)
}

func (s *clickSuite) TestInstalledNegative(c *C) {
	u, err := User()
	c.Assert(err, IsNil)
	app, err := ParseAppId("com.foo.bar_baz")
	c.Assert(err, IsNil)
	c.Check(u.Installed(app, false), Equals, false)
}

func (s *clickSuite) TestInstalledVersionNegative(c *C) {
	u, err := User()
	c.Assert(err, IsNil)
	app, err := ParseAppId("com.ubuntu.clock_clock_1000.0")
	c.Assert(err, IsNil)
	c.Check(u.Installed(app, false), Equals, false)
}

func (s *clickSuite) TestInstalledClock(c *C) {
	u, err := User()
	c.Assert(err, IsNil)
	ver := u.ccu.CGetVersion("com.ubuntu.clock")
	if ver == "" {
		c.Skip("no com.ubuntu.clock pkg installed")
	}
	app, err := ParseAppId("com.ubuntu.clock_clock")
	c.Assert(err, IsNil)
	c.Check(u.Installed(app, false), Equals, true)
	app, err = ParseAppId("com.ubuntu.clock_clock_" + ver)
	c.Assert(err, IsNil)
	c.Check(u.Installed(app, false), Equals, true)

	app, err = ParseAppId("com.ubuntu.clock_clock_10" + ver)
	c.Assert(err, IsNil)
	c.Check(u.Installed(app, false), Equals, false)

	// setVersion
	app, err = ParseAppId("com.ubuntu.clock_clock")
	c.Assert(err, IsNil)
	c.Check(u.Installed(app, true), Equals, true)
	c.Check(app.Version, Equals, ver)
}

func (s *clickSuite) TestInstalledLegacy(c *C) {
	u, err := User()
	c.Assert(err, IsNil)
	app, err := ParseAppId("_python3.4")
	c.Assert(err, IsNil)
	c.Check(u.Installed(app, false), Equals, true)
}

func (s *clickSuite) TestParseAndVerifyAppId(c *C) {
	u, err := User()
	c.Assert(err, IsNil)

	app, err := ParseAndVerifyAppId("_.foo", nil)
	c.Assert(err, Equals, ErrInvalidAppId)
	c.Check(app, IsNil)

	app, err = ParseAndVerifyAppId("com.foo.bar_baz", nil)
	c.Assert(err, IsNil)
	c.Check(app.Click, Equals, true)
	c.Check(app.Application, Equals, "baz")

	app, err = ParseAndVerifyAppId("_non-existent-app", u)
	c.Assert(err, Equals, ErrMissingAppId)
	c.Check(app, NotNil)
	c.Check(app.Original(), Equals, "_non-existent-app")

}

type helperSuite struct {
	oldHookPath        string
	symlinkPath        string
	oldHelpersDataPath string
}

var _ = Suite(&helperSuite{})

func (s *helperSuite) SetUpTest(c *C) {
	s.oldHookPath = hookPath
	hookPath = c.MkDir()
	s.symlinkPath = c.MkDir()
	s.oldHelpersDataPath = helpersDataPath
	helpersDataPath = filepath.Join(c.MkDir(), "helpers_data.json")
}

func (s *helperSuite) createHookfile(name string, content string) error {
	symlink := filepath.Join(hookPath, name) + ".json"
	filename := filepath.Join(s.symlinkPath, name)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	err = os.Symlink(filename, symlink)
	if err != nil {
		return err
	}
	return nil
}

func (s *helperSuite) createHelpersDatafile(content string) error {
	f, err := os.Create(helpersDataPath)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (s *helperSuite) TearDownTest(c *C) {
	hookPath = s.oldHookPath
	os.Remove(helpersDataPath)
	helpersDataPath = s.oldHelpersDataPath
	helpersDataMtime = time.Now().Add(-1 * time.Hour)
	helpersInfo = nil
}

func (s *helperSuite) TestHelperBasic(c *C) {
	c.Assert(s.createHelpersDatafile(`{"com.example.test_test-app_1": {"helper_id": "com.example.test_test-helper_1", "exec": "tsthlpr"}}`), IsNil)
	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "com.example.test_test-helper_1")
	c.Check(hex, Equals, "tsthlpr")
}

func (s *helperSuite) TestHelperFindsSpecific(c *C) {
	fileContent := `{"com.example.test_test-other-app": {"exec": "aaaaaaa", "helper_id": "com.example.test_aaaa-helper_1"},
    "com.example.test_test-app_1": {"exec": "tsthlpr", "helper_id": "com.example.test_test-helper_1"}}`
	c.Assert(s.createHelpersDatafile(fileContent), IsNil)

	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "com.example.test_test-helper_1")
	c.Check(hex, Equals, "tsthlpr")
}

func (s *helperSuite) TestHelperCanFail(c *C) {
	fileContent := `{"com.example.test_test-other-app": {"exec": "aaaaaaa", "helper_id": "com.example.test_aaaa-helper_1"}}`
	c.Assert(s.createHelpersDatafile(fileContent), IsNil)
	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "")
	c.Check(hex, Equals, "")
}

func (s *helperSuite) TestHelperFailInvalidJson(c *C) {
	fileContent := `{invalid json"com.example.test_test-other-app": {"exec": "aaaaaaa", "helper_id": "com.example.test_aaaa-helper_1"}}`
	c.Assert(s.createHelpersDatafile(fileContent), IsNil)
	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "")
	c.Check(hex, Equals, "")
}

func (s *clickSuite) TestHelperlegacy(c *C) {
	appname := "ubuntu-system-settings"
	app, err := ParseAppId("_" + appname)
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "")
	c.Check(hex, Equals, "")
}

// Missing Cache file test

func (s *helperSuite) TestHelperMissingCacheFile(c *C) {
	c.Assert(s.createHookfile("com.example.test_test-helper_1", `{"exec": "tsthlpr"}`), IsNil)
	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "com.example.test_test-helper_1")
	c.Check(hex, Equals, filepath.Join(s.symlinkPath, "tsthlpr"))
}

func (s *helperSuite) TestHelperFromHookBasic(c *C) {
	c.Assert(s.createHookfile("com.example.test_test-helper_1", `{"exec": "tsthlpr"}`), IsNil)
	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "com.example.test_test-helper_1")
	c.Check(hex, Equals, filepath.Join(s.symlinkPath, "tsthlpr"))
}

func (s *helperSuite) TestHelperFromHookFindsSpecific(c *C) {
	// Glob() sorts, so the first one will come first
	c.Assert(s.createHookfile("com.example.test_aaaa-helper_1", `{"exec": "aaaaaaa", "app_id": "com.example.test_test-other-app"}`), IsNil)
	c.Assert(s.createHookfile("com.example.test_test-helper_1", `{"exec": "tsthlpr", "app_id": "com.example.test_test-app"}`), IsNil)
	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "com.example.test_test-helper_1")
	c.Check(hex, Equals, filepath.Join(s.symlinkPath, "tsthlpr"))
}

func (s *helperSuite) TestHelperFromHookCanFail(c *C) {
	c.Assert(s.createHookfile("com.example.test_aaaa-helper_1", `{"exec": "aaaaaaa", "app_id": "com.example.test_test-other-app"}`), IsNil)
	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "")
	c.Check(hex, Equals, "")
}

func (s *helperSuite) TestHelperFromHookCanFailInvalidJson(c *C) {
	c.Assert(s.createHookfile("com.example.test_aaaa-helper_1", `invalid json {"exec": "aaaaaaa", "app_id": "com.example.test_test-other-app"}`), IsNil)
	app, err := ParseAppId("com.example.test_test-app_1")
	c.Assert(err, IsNil)
	hid, hex := app.Helper()
	c.Check(hid, Equals, "")
	c.Check(hex, Equals, "")
}
