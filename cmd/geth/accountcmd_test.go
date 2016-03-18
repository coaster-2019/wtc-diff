// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cespare/cp"
)

// These tests are 'smoke tests' for the account related
// subcommands and flags.
//
// For most tests, the test files from package accounts
// are copied into a temporary keystore directory.

func tmpDatadirWithKeystore(t *testing.T) string {
	datadir := tmpdir(t)
	keystore := filepath.Join(datadir, "keystore")
	source := filepath.Join("..", "..", "accounts", "testdata", "keystore")
	if err := cp.CopyAll(keystore, source); err != nil {
		t.Fatal(err)
	}
	return datadir
}

func TestAccountListEmpty(t *testing.T) {
	geth := runGeth(t, "account")
	geth.expectExit()
}

func TestAccountList(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t, "--datadir", datadir, "account")
	defer geth.expectExit()
	geth.expect(`
Account #0: {7ef5a6135f1fd6a02593eedc869c6d41d934aef8}
Account #1: {f466859ead1932d743d622cb74fc058882e8648a}
Account #2: {289d485d9771714cce91d3393d764e1311907acc}
`)
}

func TestAccountNew(t *testing.T) {
	geth := runGeth(t, "--lightkdf", "account", "new")
	defer geth.expectExit()
	geth.expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Repeat passphrase: {{.InputLine "foobar"}}
`)
	geth.expectRegexp(`Address: \{[0-9a-f]{40}\}\n`)
}

func TestAccountNewBadRepeat(t *testing.T) {
	geth := runGeth(t, "--lightkdf", "account", "new")
	defer geth.expectExit()
	geth.expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "something"}}
Repeat passphrase: {{.InputLine "something else"}}
Fatal: Passphrases do not match
`)
}

func TestAccountUpdate(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t,
		"--datadir", datadir, "--lightkdf",
		"account", "update", "f466859ead1932d743d622cb74fc058882e8648a")
	defer geth.expectExit()
	geth.expect(`
Unlocking account f466859ead1932d743d622cb74fc058882e8648a | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Please give a new password. Do not forget this password.
Passphrase: {{.InputLine "foobar2"}}
Repeat passphrase: {{.InputLine "foobar2"}}
`)
}

func TestWalletImport(t *testing.T) {
	geth := runGeth(t, "--lightkdf", "wallet", "import", "testdata/guswallet.json")
	defer geth.expectExit()
	geth.expect(`
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foo"}}
Address: {d4584b5f6229b7be90727b0fc8c6b91bb427821f}
`)

	files, err := ioutil.ReadDir(filepath.Join(geth.Datadir, "keystore"))
	if len(files) != 1 {
		t.Errorf("expected one key file in keystore directory, found %d files (error: %v)", len(files), err)
	}
}

func TestWalletImportBadPassword(t *testing.T) {
	geth := runGeth(t, "--lightkdf", "wallet", "import", "testdata/guswallet.json")
	defer geth.expectExit()
	geth.expect(`
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "wrong"}}
Fatal: Could not create the account: Decryption failed: PKCS7Unpad failed after AES decryption
`)
}

func TestUnlockFlag(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--dev",
		"--unlock", "f466859ead1932d743d622cb74fc058882e8648a",
		"js", "testdata/empty.js")
	geth.expect(`
Unlocking account f466859ead1932d743d622cb74fc058882e8648a | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
`)
	geth.expectExit()

	wantMessages := []string{
		"Unlocked account f466859ead1932d743d622cb74fc058882e8648a",
	}
	for _, m := range wantMessages {
		if strings.Index(geth.stderrText(), m) == -1 {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagWrongPassword(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--dev",
		"--unlock", "f466859ead1932d743d622cb74fc058882e8648a")
	defer geth.expectExit()
	geth.expect(`
Unlocking account f466859ead1932d743d622cb74fc058882e8648a | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "wrong1"}}
Unlocking account f466859ead1932d743d622cb74fc058882e8648a | Attempt 2/3
Passphrase: {{.InputLine "wrong2"}}
Unlocking account f466859ead1932d743d622cb74fc058882e8648a | Attempt 3/3
Passphrase: {{.InputLine "wrong3"}}
Fatal: Failed to unlock account: f466859ead1932d743d622cb74fc058882e8648a
`)
}

// https://github.com/ethereum/go-ethereum/issues/1785
func TestUnlockFlagMultiIndex(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--dev",
		"--unlock", "0,2",
		"js", "testdata/empty.js")
	geth.expect(`
Unlocking account 0 | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Unlocking account 2 | Attempt 1/3
Passphrase: {{.InputLine "foobar"}}
`)
	geth.expectExit()

	wantMessages := []string{
		"Unlocked account 7ef5a6135f1fd6a02593eedc869c6d41d934aef8",
		"Unlocked account 289d485d9771714cce91d3393d764e1311907acc",
	}
	for _, m := range wantMessages {
		if strings.Index(geth.stderrText(), m) == -1 {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagPasswordFile(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--dev",
		"--password", "testdata/passwords.txt", "--unlock", "0,2",
		"js", "testdata/empty.js")
	geth.expectExit()

	wantMessages := []string{
		"Unlocked account 7ef5a6135f1fd6a02593eedc869c6d41d934aef8",
		"Unlocked account 289d485d9771714cce91d3393d764e1311907acc",
	}
	for _, m := range wantMessages {
		if strings.Index(geth.stderrText(), m) == -1 {
			t.Errorf("stderr text does not contain %q", m)
		}
	}
}

func TestUnlockFlagPasswordFileWrongPassword(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t,
		"--datadir", datadir, "--nat", "none", "--nodiscover", "--dev",
		"--password", "testdata/wrong-passwords.txt", "--unlock", "0,2")
	defer geth.expectExit()
	geth.expect(`
Fatal: Failed to unlock account: 0
`)
}
