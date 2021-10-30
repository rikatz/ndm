package photon

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	defaultPhotonURL  = "https://packages.vmware.com/photon/4.0/Rev1/ova/photon-ova-4.0-ca7c9e9330.ova"
	defaultPhotonSHA1 = "cb356cc14435e826d3db82aebe5eda2ee3d731d2"
)

func CheckPhotonOVA(url string) (string, error) {
	downloadURL := defaultPhotonURL
	if url != "" {
		downloadURL = url
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := fmt.Sprintf("%s/.cache/ndm", homedir)
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%s/photon.ova", cacheDir)

	_, err = os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		log.Printf("Downloading %s into %s", downloadURL, filename)
		err := downloadOVA(filename, downloadURL)
		if err != nil {
			return "", err
		}
	}
	log.Printf("Verifying SHA1 of file %s", filename)
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	if hex.EncodeToString(h.Sum(nil)) != defaultPhotonSHA1 {
		return "", fmt.Errorf("SHA1 of OVA file %s is invalid, please remove the file and re-run the program", filename)
	}
	log.Printf("File %s downloaded and checksum %s matches", filename, defaultPhotonSHA1)
	return filename, nil
}

func downloadOVA(file, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode > 305 {
		return fmt.Errorf("could not download the Photon file %s", url)
	}
	defer resp.Body.Close()

	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
