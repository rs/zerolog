// +build !windows

package journald_test

import "github.com/rs/zerolog"
import "github.com/rs/zerolog/journald"

func ExampleJournalDWriter () {
	log := zerolog.New(journald.NewJournalDWriter())
	log.Info().Str("foo", "bar").Msg("Journal Test")
	// Output: 
}

/*

There is no automated way to verify the output - since the output is sent
to journald process and method to retrieve is journalctl. Will find a way
to automate the process and fix this test 

$ journalctl -a -o verbose -f

Mon 2018-04-23 07:12:26.635654 PDT [s=349577f78be941179a3aa5d849b2cefa;i=cdc;b=30c80d57c6b7443fb0481ba66e9b7a62;m=9412558ee4;t=56a84a0a03a5f;x=2ba55bfcdc33820f]
    PRIORITY=6
    _BOOT_ID=30c80d57c6b7443fb0481ba66e9b7a62
    _MACHINE_ID=4acf281d5c48d411a26b176659c09a09
    _HOSTNAME=minion1
    _CAP_EFFECTIVE=0
    _AUDIT_LOGINUID=1000
    _TRANSPORT=journal
    _GID=1000
    _UID=1000
    _AUDIT_SESSION=622
    MESSAGE=Journal Test
    FOO=bar
    JSON={"level":"info","foo":"bar","message":"Journal Test"}
    _COMM=j2
    _PID=21466
    _SOURCE_REALTIME_TIMESTAMP=1524492746635654
*/

