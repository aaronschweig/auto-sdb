package helpers

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if !encountered[elements[v]] {

			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}
