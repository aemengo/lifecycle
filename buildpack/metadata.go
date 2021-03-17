package buildpack

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/buildpacks/lifecycle/api"
)

const (
	warningTypesKey      = "Warning: types table isn't supported in this buildpack api version. The launch, build and cache flags should be in the top level. Ignoring the values in the types table."
	warningTypesTopLevel = "the launch, cache and build flags should be in the types table of %s"
)

type LayerMetadataFile struct {
	Data   interface{} `json:"data" toml:"metadata"`
	Build  bool        `json:"build" toml:"build"`
	Launch bool        `json:"launch" toml:"launch"`
	Cache  bool        `json:"cache" toml:"cache"`
}

func EncodeLayerMetadataFile(path string, meta *LayerMetadataFile, buildpackAPI string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if api.MustParse(buildpackAPI).Compare(api.MustParse("0.6")) < 0 {
		// v05
		return toml.NewEncoder(f).Encode(meta)
	}

	// v06
	return toml.NewEncoder(f).Encode(struct {
		Data  interface{} `json:"data" toml:"metadata"`
		Types struct {
			Build  bool `json:"build" toml:"build"`
			Launch bool `json:"launch" toml:"launch"`
			Cache  bool `json:"cache" toml:"cache"`
		} `json:"types" toml:"types"`
	}{
		Data: meta.Data,
	})
}

// TODO: pass the logger and print the warning inside (instead of returning a message)
// Anthony: since we're gonna pass a logger, maybe we should create an struct
//          and instantiate with buildpackAPI too. I find ', buildpackAPI string)'
//          part to be ruining the signature.
func DecodeLayerMetadataFile(path string, dest interface{}, buildpackAPI string) (string, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	defer f.Close()

	if api.MustParse(buildpackAPI).Compare(api.MustParse("0.6")) < 0 {
		// v05
		return decodeLayerMetadataFileO5(path, dest)
	}

	// v06
	return decodeLayerMetadataFileO6(path, dest)
}

func decodeLayerMetadataFileO5(path string, dest interface{}) (string, error) {
	md, err := toml.DecodeFile(path, dest)
	if err != nil {
		return "", err
	}

	if md.IsDefined("types") {
		return warningTypesKey, nil
	}

	return "", nil
}

func decodeLayerMetadataFileO6(path string, dest interface{}) (string, error) {
	var metaWithTypes struct {
		Data  interface{} `json:"data" toml:"metadata"`
		Types struct {
			Build  bool `json:"build" toml:"build"`
			Launch bool `json:"launch" toml:"launch"`
			Cache  bool `json:"cache" toml:"cache"`
		} `json:"types" toml:"types"`
	}

	md, err := toml.DecodeFile(path, &metaWithTypes)
	if err != nil {
		return "", err
	}

	if md.IsDefined("build") || md.IsDefined("launch") || md.IsDefined("cache") {
		return fmt.Sprintf(warningTypesTopLevel, path), nil
	}

	meta, ok := dest.(*LayerMetadataFile)
	if !ok {
		//TODO: handle better
		//but this would be a programmer error
		panic("you passed the wrong type of file in")
	}

	meta.Data = metaWithTypes.Data
	meta.Build = metaWithTypes.Types.Build
	meta.Launch = metaWithTypes.Types.Launch
	meta.Cache = metaWithTypes.Types.Cache
	return "", nil
}
