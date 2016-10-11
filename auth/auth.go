package auth

import "os"

// WriteToken writes an auth token to a local file.
func WriteToken(token, filePath string) error {
	tFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	tFile.Chmod(0600)
	defer tFile.Close()

	tFile.WriteString(token)

	return nil
}
