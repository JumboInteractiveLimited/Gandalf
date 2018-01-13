// One Tool to rule them all, One Tool to CI them,
// One Tool to test them all and in the darkness +1 them.
//
// Gandalf is designed to provide a language and stack agnostic HTTP API contract
// testing suite and prototyping toolchain. This is achieved by; running an HTTP
// API (aka provider), connecting to it as a real client (aka consumer) of the
// provider, asserting that it matches various rules (aka contracts). Optionally,
// once a contract is written you can then generate an approximation of the API
// (this happens just before the contract is tested) in the form of a mock. This
// allows for rapid prototyping and/or parallel development of the real consumer
// and provider implementations.
//
// Gandalf has no allegiance to any specific paradigms, technologies, or concepts
// and should bend to fit real world use cases as opposed to vice versa. This
// means if Gandalf does something one way today it does not mean that tomorrow it
// could not support a different way provided someone has a use for it.
//
// While Gandalf does use golang and the go test framework, it is not specific to
// go as at its core it just makes HTTP requests and checks the responses. Your
// web server or clients can be written in any language/framework. The official
// documentation also uses JSON and RESTful API's as examples but Gandalf supports
// any and all paradigms or styles of API.
//
// Most go programs are compiled down to a binary and executed, Gandalf is designed
// to be used as a library to write your own tests and decorate the test binary
// instead. For example, Gandalf does have several command line switches however
// they are provided to the `go test` command instead of some non existent `Gandalf`
// command. This allows Gandalf to get all kind of testing and benchmarking support
// for free while being a well known stable base to build upon.
//
// Contract testing can be a bit nebulous and also has various option prefixes such
// as Consumer Driven, Gandalf cares not for any prefixes (who writes contracts and
// where is up to you) nor does it care if you are testing the interface or your
// API or the business logic or some combination of both, no one will save you from
// blowing your own foot off if you choose to.
package gandalf

import (
	"fmt"
	"os"
	"testing"

	"github.com/fatih/color"
)

// Main should be run by TestMain in order for Gandalf to
// analyze the whole test run.
//  func TestMain(m *testing.M) {
//    gandalf.Main(m)
//  }
func Main(m *testing.M) {
	displayMascot()
	if OverrideColour {
		color.NoColor = false
	}
	ret := m.Run()
	if ret == 0 {
		color.Green("✔ He that breaks a thing to find out what it is has left the path of wisdom.")
	} else {
		color.Red("✘ YOU SHALL NOT PASS!")
	}
	os.Exit(ret)
}

func displayMascot() {
	fmt.Println(`
          -~~-                                                ,#@
                                      ,#5ppppppppS##Ww,,,,s###S##
          -~~ .                     sTpSZ##############S#######@
       '  #s,         .            #pS@##############@#######Q#
          Q#Q}@o  \               #pS@###########SQ###***"7^^
     [   ,p@@#@#~   [  [          bS@#############SkSb
     '   @####S@  /    ^         ]b#################8@
      '   %@@@M                  ####################SQ                      #S[
  \   ^.  ~ 8#b  '              #S####################Sk,                  {#SS[
             #@               ,#S###############@QQSSSSSQ@@w              ,|#SS
        -    @#@            ,#b######QS669SSps#SS##########SS##Wg         #|SS^
         '%mwJQ@       ,,s#@Qs###8bppSSSSSZ######$QQQQ$##########S#W,    #|SS"
            ^*@QppppppppppSSSSSSSSSSSSSSZ@Q#M*7^^  ,-~w ~^Qj%WQQ####S#w {|SSf
               #SC"7*@m@QQQQQQQQ#k==%7"=o    - ^       *   3    [^"%@QS#QSS#
               1#@   j  |@   p    aMw    ~  ~    #@m    ~   "  ,^     |5#M#
                QSQ   '*g    \     @a   ,   \      #   /     {^       ##SS
                1@#     L     ^.      ,' s @@Q~.    .~               $|SS
                 @#b    Q              ob   9b".             ab     ;|SS
                 jS#   ] ^w         ,a2, ,Z>o. ,7*-,,  ,.=^.  }     #SS
                  @#b  [    *^"^^"b*    b  > ]^               ]    ##S"
                  ^#8  t                ^~^\-^                ]   6#SL
                   @#b |   o                               j  @  ]|S#
                    ##  ~  Q    @                    [     {  @  |S#
                    @#k 1  %    ]                    k     b C]###Q
                     ##w#1, Q    p                  .b    f.C jyR#1%@m,
                    #@#@{  \^Q   1                  f    @4    @#v@
                     *#S#    "Q   @                ]~   A      #^^
                      1##      \   V              )b  ;\       b
                       Qb       'V  %            ;^ ,^         b
                       1b          "Q^m         4,*\           b
                        >              |w     A                b
                        .                 "Q C                 ~
                        b                                     ]
                        1                                     b
                         t                                   C
                          V                                 C
                           ^w                            ,#,
                          #b ."w,                     .4b. ,k
                         1Z~ Z^    ^"~--.......-~<^      \,jm

	`)
}
