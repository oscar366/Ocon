package updating

import "net/http"
import "encoding/json"//for github api
import "fmt"
import "io"
import "strconv"
import "os"
import "time"

var necesitoactualizar = false

//version
func fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
//reinstall the file for updates
func ReinstallOcon() {
	if !necesitoactualizar {
		fmt.Println("tú no necesitas actualizar (you don't need an update) 🐟")
		return
	}

	// Get the latest release
	data, err := fetch("https://api.github.com/repos/oscar366/ocon/releases/latest")
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return
	}

	// Get assets_url
	var release struct {
		AssetsURL string `json:"assets_url"`
	}

	if err := json.Unmarshal(data, &release); err != nil {
		fmt.Println("Error parsing release:", err)
		return
	}

	// Get the release assets
	data2, err := fetch(release.AssetsURL)
	if err != nil {
		fmt.Println("Error getting release assets:", err)
		return
	}

	// Parse assets
	var assets []ReleaseAsset

	if err := json.Unmarshal(data2, &assets); err != nil {
		fmt.Println("Error parsing assets:", err)
		return
	}

	if len(assets) == 0 {
		fmt.Println("No release assets found.")
		return
	}

	// Display assets
	fmt.Println("Which one do I install? (type number):")

	for i, asset := range assets {
		fmt.Printf("%d: %s\n", i, asset.Name)
	}

	var userInput string
	fmt.Scan(&userInput)

	num, err := strconv.Atoi(userInput)
	if err != nil {
		fmt.Println("Invalid number:", err)
		return
	}

	if num < 0 || num >= len(assets) {
		fmt.Println("Out of bounds")
		return
	}

	// Download selected asset
	data, err = fetch(assets[num].BrowserDownloadURL)
	if err != nil {
		fmt.Println("Download failed:", err)
		return
	}

	err = os.WriteFile("ocon", data, 0755)
	if err != nil {
		fmt.Println("Failed to write file:", err)
		return
	}

	fmt.Println("==================")
	fmt.Println("Steps After Update")
	fmt.Println("==================")
	fmt.Println("Ocon has downloaded the new version to the same location as the older version.")
	fmt.Println("1. Find the newly downloaded file (called ocon with no extension).")
	fmt.Println("2. Add the appropriate file extension for your operating system.")
	fmt.Println("   - Windows: .exe")
	fmt.Println("   - Linux: no extension required")
	fmt.Println("3. Delete the older version of Ocon.")
	fmt.Println("4. Run the new version.")
}


type Update struct {
	TagName string `json:"tag_name"`
	Name string `json:"name"`
	Body string `json:body`
}
//if needs update
func IfNecesitoActualizar(Ver string) {
	url := "https://api.github.com/repos/oscar366/Ocon/releases/latest"
	
	if time.Now().Second() % 2 != 0 {
		return
	}
	
	data, err := fetch(url)
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return
	}

	var u Update

	if err := json.Unmarshal(data, &u); err != nil {
		fmt.Println("Error parsing update JSON:", err)
		return
	}

	if u.TagName != Ver {
		fmt.Println("\033[31mYou need an update.\033[0m")
		fmt.Println("Current:", Ver)
		fmt.Println("Latest:", u.TagName)
		necesitoactualizar = true
	}
}