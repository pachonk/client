package engine

import (
	"strings"
	"testing"

	"github.com/keybase/client/go/libkb"
)

func TestPGPEncrypt(t *testing.T) {
	tc := SetupEngineTest(t, "PGPEncrypt")
	defer tc.Cleanup()

	u := createFakeUserWithPGPSibkey(tc)
	trackUI := &FakeIdentifyUI{
		Proofs: make(map[string]string),
	}
	ctx := &Context{IdentifyUI: trackUI, SecretUI: u.NewSecretUI()}

	sink := libkb.NewBufferCloser()
	arg := &PGPEncryptArg{
		Recips: []string{"t_alice", "t_bob+kbtester1@twitter", "t_charlie+tacovontaco@twitter"},
		Source: strings.NewReader("track and encrypt, track and encrypt"),
		Sink:   sink,
		NoSign: true,
	}

	eng := NewPGPEncrypt(arg, tc.G)
	if err := RunEngine(eng, ctx); err != nil {
		t.Fatal(err)
	}

	out := sink.Bytes()
	if len(out) == 0 {
		t.Fatal("no output")
	}

	if err := runUntrack(tc.G, u, "t_alice"); err != nil {
		t.Fatal(err)
	}
	if err := runUntrack(tc.G, u, "t_bob"); err != nil {
		t.Fatal(err)
	}
	if err := runUntrack(tc.G, u, "t_charlie"); err != nil {
		t.Fatal(err)
	}
}

func TestPGPEncryptSelfNoKey(t *testing.T) {
	tc := SetupEngineTest(t, "PGPEncrypt")
	defer tc.Cleanup()

	u := CreateAndSignupFakeUser(tc, "login")
	trackUI := &FakeIdentifyUI{
		Proofs: make(map[string]string),
	}
	ctx := &Context{IdentifyUI: trackUI, SecretUI: u.NewSecretUI()}

	sink := libkb.NewBufferCloser()
	arg := &PGPEncryptArg{
		Recips: []string{"t_alice", "t_bob+kbtester1@twitter", "t_charlie+tacovontaco@twitter"},
		Source: strings.NewReader("track and encrypt, track and encrypt"),
		Sink:   sink,
		NoSign: true,
	}

	eng := NewPGPEncrypt(arg, tc.G)
	err := RunEngine(eng, ctx)
	if err == nil {
		t.Fatal("no error encrypting for self without pgp key")
	}
	if _, ok := err.(libkb.NoKeyError); !ok {
		t.Fatalf("expected error type libkb.NoKeyError, got %T (%s)", err, err)
	}
}

func TestPGPEncryptNoTrack(t *testing.T) {
	tc := SetupEngineTest(t, "PGPEncrypt")
	defer tc.Cleanup()

	u := createFakeUserWithPGPSibkey(tc)
	trackUI := &FakeIdentifyUI{
		Proofs: make(map[string]string),
	}
	ctx := &Context{IdentifyUI: trackUI, SecretUI: u.NewSecretUI()}

	sink := libkb.NewBufferCloser()
	arg := &PGPEncryptArg{
		Recips:    []string{"t_alice", "t_bob+kbtester1@twitter", "t_charlie+tacovontaco@twitter"},
		Source:    strings.NewReader("identify and encrypt, identify and encrypt"),
		Sink:      sink,
		NoSign:    true,
		SkipTrack: true,
	}

	eng := NewPGPEncrypt(arg, tc.G)
	if err := RunEngine(eng, ctx); err != nil {
		t.Fatal(err)
	}

	out := sink.Bytes()
	if len(out) == 0 {
		t.Fatal("no output")
	}

	assertNotTracking(t, "t_alice")
	assertNotTracking(t, "t_bob")
	assertNotTracking(t, "t_charlie")
}
