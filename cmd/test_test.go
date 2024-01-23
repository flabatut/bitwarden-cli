package cmd_test

// TODO

// https://g14a.dev/posts/Testing-Cobra-Subcommands/
// https://stackoverflow.com/questions/59709345/how-to-implement-unit-tests-for-cli-commands-in-go/59714127#59714127
// https://github.com/helm/helm/tree/main/cmd/helm
// https://github.com/spf13/cobra/issues/770
// https://nayaktapan37.medium.com/testing-cobra-commands-in-golang-ca1fe4ad6657
// https://gianarb.it/blog/golang-mockmania-cli-command-with-cobra
// https://clavinjune.dev/en/blogs/implement-unit-test-for-cli-apps-using-golang-and-cobra/
// https://www.bradcypert.com/testing-a-cobra-cli-in-go/

// import (
// 	"bytes"
// 	"fmt"
// 	"testing"

// 	"github.com/spf13/cobra"
// )

// func emptyRun(*cobra.Command, []string) {}

// func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
// 	_, output, err = executeCommandC(root, args...)
// 	return output, err
// }

// func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
// 	buf := new(bytes.Buffer)
// 	root.SetOut(buf)
// 	root.SetErr(buf)
// 	root.SetArgs(args)

// 	c, err = root.ExecuteC()

// 	return c, buf.String(), err
// }

// func NewRootCmd(in string) *cobra.Command {
//     return &cobra.Command{
//       Use:   "hugo",
//       Short: "Hugo is a very fast static site generator",
//       Long: `A Fast and Flexible Static Site Generator built with
//                 love by spf13 and friends in Go.
//                 Complete documentation is available at http://hugo.spf13.com`,
//       RunE: func(cmd *cobra.Command, args []string) (error) {
//           fmt.Fprintf(cmd.OutOrStdout(), in)
//           return nil
//       },
// 	}
// }

// func TestTestCommand(t *testing.T) {
// 	fmt.Println("echo")

// 	// var rootCmdArgs []string
// 	// rootCmd := &cobra.Command{
// 	// 	Use:  "test",
// 	// 	Args: cobra.ExactArgs(0),
// 	// 	Run:  func(_ *cobra.Command, args []string) { rootCmdArgs = args },
// 	// }
// 	// aCmd := &cobra.Command{Use: "a", Args: cobra.NoArgs, Run: emptyRun}
// 	// bCmd := &cobra.Command{Use: "b", Args: cobra.NoArgs, Run: emptyRun}
// 	// rootCmd.AddCommand(aCmd, bCmd)

// 	// output, err := executeCommand(rootCmd, "one", "two")
// 	// if output != "" {
// 	// 	t.Errorf("Unexpected output: %v", output)
// 	// }
// 	// if err != nil {
// 	// 	t.Errorf("Unexpected error: %v", err)
// 	// }

// 	// got := strings.Join(rootCmdArgs, " ")
// 	// if got != onetwo {
// 	// 	t.Errorf("rootCmdArgs expected: %q, got: %q", onetwo, got)
// 	// }
// }
