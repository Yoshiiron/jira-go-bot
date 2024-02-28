package HelpingFuncs

//func main() {
//	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
//	integers := "1234567890"
//
//}

//func random() string {
//	rand.New(rand.NewSource(time.Now().UnixNano()))
//	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
//		"abcdefghijklmnopqrstuvwxyz" +
//		"0123456789")
//	length := 8
//	var b strings.Builder
//	for i := 0; i < length; i++ {
//		b.WriteRune(chars[rand.Intn(len(chars))])
//	}
//	str := b.String() // Например "ExcbsVQs"
//	return str
//}

//func randomCharsFromString() string {
//	a, err := randomizer.RandomString(6)
//	if err != nil {
//		fmt.Println(1)
//	}
//	return a
//
//}

func ReturnKeys(m map[string]string) string {

	keys := make([]string, 0, len(m))
	var keyreturner string

	for key := range m {
		keys = append(keys, key)
	}

	for _, key := range keys {
		keyreturner = key
	}

	return keyreturner
}
