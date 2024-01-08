// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package sketch

import (
	"context"
	"fmt"
	"os"

	"github.com/arduino/arduino-cli/commands/sketch"
	"github.com/arduino/arduino-cli/internal/cli/arguments"
	"github.com/arduino/arduino-cli/internal/cli/feedback"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initArchiveCommand creates a new `archive` command
func initArchiveCommand() *cobra.Command {
	var includeBuildDir, overwrite bool

	archiveCommand := &cobra.Command{
		Use:   fmt.Sprintf("archive <%s> <%s>", tr("sketchPath"), tr("archivePath")),
		Short: tr("Creates a zip file containing all sketch files."),
		Long:  tr("Creates a zip file containing all sketch files."),
		Example: "" +
			"  " + os.Args[0] + " archive\n" +
			"  " + os.Args[0] + " archive .\n" +
			"  " + os.Args[0] + " archive . MySketchArchive.zip\n" +
			"  " + os.Args[0] + " archive /home/user/Arduino/MySketch\n" +
			"  " + os.Args[0] + " archive /home/user/Arduino/MySketch /home/user/MySketchArchive.zip",
		Args: cobra.MaximumNArgs(2),
		Run:  func(cmd *cobra.Command, args []string) { runArchiveCommand(args, includeBuildDir, overwrite) },
	}

	archiveCommand.Flags().BoolVar(&includeBuildDir, "include-build-dir", false, tr("Includes %s directory in the archive.", "build"))
	archiveCommand.Flags().BoolVarP(&overwrite, "overwrite", "f", false, tr("Overwrites an already existing archive"))

	return archiveCommand
}

func runArchiveCommand(args []string, includeBuildDir bool, overwrite bool) {
	logrus.Info("Executing `arduino-cli sketch archive`")

	sketchPathArg := ""
	if len(args) > 0 {
		sketchPathArg = args[0]
	}

	archivePathArg := ""
	if len(args) > 1 {
		archivePathArg = args[1]
	}

	sketchPath := arguments.InitSketchPath(sketchPathArg)
	sk, err := sketch.LoadSketch(context.Background(), &rpc.LoadSketchRequest{SketchPath: sketchPath.String()})
	if err != nil {
		feedback.FatalError(err, feedback.ErrGeneric)
	}
	feedback.WarnAboutDeprecatedFiles(sk)

	if _, err := sketch.ArchiveSketch(context.Background(),
		&rpc.ArchiveSketchRequest{
			SketchPath:      sketchPath.String(),
			ArchivePath:     archivePathArg,
			IncludeBuildDir: includeBuildDir,
			Overwrite:       overwrite,
		},
	); err != nil {
		feedback.Fatal(tr("Error archiving: %v", err), feedback.ErrGeneric)
	}
}
