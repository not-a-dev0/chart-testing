// Copyright The Helm Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/helm/chart-testing/pkg/chart"
	"github.com/helm/chart-testing/pkg/config"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install and test a chart",
		Long: heredoc.Doc(`
			Run 'helm install' and ' helm test' on

			* changed charts (default)
			* specific charts (--charts)
			* all charts (--all)

			in given chart directories.

			Charts may have multiple custom values files matching the glob pattern
			'*-values.yaml' in a directory named 'ci' in the root of the chart's
			directory. The chart is installed and tested for each of these files.
			If no custom values file is present, the chart is installed and
			tested with defaults.`),
		Run: install,
	}

	flags := cmd.Flags()
	addInstallFlags(flags)
	addCommonLintAndInstallFlags(flags)
	return cmd
}

func addInstallFlags(flags *flag.FlagSet) {
	flags.String("build-id", "", heredoc.Doc(`
		An optional, arbitrary identifier that is added to the name of the namespace a
		chart is installed into. In a CI environment, this could be the build number or
		the ID of a pull request. If not specified, the name of the chart is used`))
	flags.String("helm-extra-args", "", heredoc.Doc(`
		Additional arguments for Helm. Must be passed as a single quoted string
		(e.g. "--timeout 500 --tiller-namespace tiller"`))
	flags.String("namespace", "", heredoc.Doc(`
		Namespace to install the release(s) into. If not specified, each release will be
		installed in its own randomly generated namespace.`))
	flags.String("release-label", "app.kubernetes.io/instance", heredoc.Doc(`
		The label to be used as a selector when inspecting resources created by charts.
		This is only used if namespace is specified.`))
}

func install(cmd *cobra.Command, args []string) {
	fmt.Println("Installing charts...")

	configuration, err := config.LoadConfiguration(cfgFile, cmd)
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		os.Exit(1)
	}

	testing := chart.NewTesting(*configuration)
	results, err := testing.InstallCharts()
	if err != nil {
		fmt.Printf("Error installing charts: %s\n", err)
	} else {
		fmt.Println("All charts installed successfully")
	}

	testing.PrintResults(results)

	if err != nil {
		os.Exit(1)
	}
}
