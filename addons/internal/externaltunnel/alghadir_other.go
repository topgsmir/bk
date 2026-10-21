//go:build !linux

package externaltunnel

import (
	"context"
	"fmt"
)

func udp2rawAsset(arch string) (Asset, error) { return Asset{}, fmt.Errorf("Alghadir requires Linux") }
func runAlghadir(ctx context.Context, s Spec, dir string) error {
	return fmt.Errorf("Alghadir requires Linux")
}
func AlghadirChild(ctx context.Context, s Spec, ipc string) error {
	return fmt.Errorf("Alghadir requires Linux")
}
