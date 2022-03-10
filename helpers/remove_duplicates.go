package helpers

func RemoveDuplicates(elements []string) []string {
	encountered := make(map[string]struct{}, len(elements))
	result := make([]string, 0, len(elements))

	for v := range elements {
		if _, ok := encountered[elements[v]]; !ok {
			encountered[elements[v]] = struct{}{}
			result = append(result, elements[v])
		}
	}
	return result
}
