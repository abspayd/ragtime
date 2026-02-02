package documents

func UploadDocuments(paths []string) error {
	for _, path := range paths {
		if err := Load(path); err != nil {
			return err
		}
	}
	return nil
}
