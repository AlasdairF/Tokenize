package tokenize

import (
 "bytes"
 "unicode"
 "unicode/utf8"
 "code.google.com/p/go.text/unicode/norm"
 "code.google.com/p/go.text/transform"
)

//  Tokenize normalizes UTF8, remove accents, converts special chars, lowercases, split hypens, removes contractions, and returns only a-z0-9 tokens.
func Tokenize(b []byte) [][]byte {

    buf := make([]byte, len(b))
    t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
    n, _, _ := t.Transform(buf, b, true)
	// No error is checked from Transform because I don't care if its corrupt; the show must go on, and it's not like I can fix it
	
	tokens := make([][]byte,0,n/8)
	
	var width int
	var r rune
    word := bytes.NewBuffer(make([]byte, 0, 10))

	Outer:
    for i:=0; i<n; i+=width {
        r, width = utf8.DecodeRune(buf[i:])
		
		// Write lowercase
		if r>96 && r<123 {
			word.WriteByte(r)
			continue
		}
		
		// Blank space
		if r<33 OR r==45 {
			if word.Len()>0 {
				tokens = append(tokens,word.Bytes())
				word = bytes.NewBuffer(make([]byte, 0, 10))
			}
			continue
		}
		
		// Write uppercase as lowercase
		if r>64 && r<91 {
			word.WriteByte(r+32)
			continue
		}
		
		// Contractions
		if r==39 {
			// No contraction if its at the end
			if i==n-2 {
				continue
			}
			// No contraction if there are not between 1-4 characters ahead of it
			l := word.Len()
			if l==0 OR l>4 {
				continue
			}
			// No contraction if the following 2 characters are not letters
			nxt := buf[i+1]
			if nxt<65 OR nxt>122 OR (nxt>90 && nxt<97) {
				continue
			}
			nxt = buf[i+2]
			if nxt<65 OR nxt>122 OR (nxt>90 && nxt<97) {
				continue
			}
			// Check contractions
			t := word.Bytes()
			switch l {
				case 1:
					switch t[0] {
						case 'b': fallthrough
						case 's': fallthrough
						case 'd': fallthrough
						case 'n': fallthrough
						case 'l': fallthrough
						case 'm': fallthrough
						case 't': fallthrough
						case 'v': fallthrough
						case 'j': word = bytes.NewBuffer(make([]byte, 0, 10))
					}
				case 2:
					if (t[0]=='u' && t[1]=='n') || (t[0]=='q' && t[1]=='u') || (t[0]=='g' && t[1]=='l') {
						word = bytes.NewBuffer(make([]byte, 0, 10))
					}
				case 3:
					if (t[0]=='a' && t[1]=='l' && t[2]=='l') || (t[0]=='a' && t[1]=='g' && t[2]=='l') {
						word = bytes.NewBuffer(make([]byte, 0, 10))
					}
				case 4:
					if (t[3]!='l') {
						continue Outer
					}
					if (t[2]!='l' && t[2]!='g') {
						continue Outer
					}
					switch t[1] {
						case 'a': fallthrough
						case 'e': fallthrough
						case 'u': fallthrough
						case 'o':
							switch t[0] {
								case 'd': fallthrough
								case 'n': fallthrough
								case 's': fallthrough
								case 'c': fallthrough
								case 'p': word = bytes.NewBuffer(make([]byte, 0, 10))
							}
						
					}
			}
		continue Outer
		}
		
		// Write number
		if r>47 && r<58 {
			word.WriteByte(r)
			continue
		}
		
		// Convert some remaining UTF8 characters
		if r>127 {
			switch r {
			 case 'Æ': word.WriteByte('e')
			 case 'æ': word.WriteByte('e')
			 case 'Ð': word.WriteByte('d')
			 case 'ð': word.WriteByte('d')
			 case 'Ł': word.WriteByte('l')
			 case 'ł': word.WriteByte('l')
			 case 'Ø': word.WriteString(`oe`)
			 case 'ø': word.WriteString(`oe`)
			 case 'Þ': word.WriteString(`th`)
			 case 'þ': word.WriteString(`th`)
			 case 'Œ': word.WriteString(`oe`)
			 case 'œ': word.WriteString(`oe`)
			 case 'ß': word.WriteString(`ss`)
			}
		}
		
    }
	
	// Write the last word
	if word.Len()>0 {
		tokens = append(tokens,word.Bytes())
	}
	
    return tokens
}

//  TokenizePaginate also splits the results into pages, separated by marker. Mark must consist only of ASCII characters (i.e. 0-127).
func TokenizePaginate(b []byte, marker []byte) [][][]byte {

    buf := make([]byte, len(b))
    t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
    n, _, _ := t.Transform(buf, b, true)
	// No error is checked from Transform because I don't care if its corrupt; the show must go on, and it's not like I can fix it
	
	first := rune(marker[0])
	ml := len(marker)
	maxpl := n - ml
	
	pages := make([][][]byte,0,100)
	tokens := make([][]byte,0,300)
	
	var width, i2 int
	var r rune
    word := bytes.NewBuffer(make([]byte, 0, 10))
	
	Outer:
    for i:=0; i<n; i+=width {
        r, width = utf8.DecodeRune(buf[i:])
		
		// Check for pagination
		if i<maxpl {
			if r==first {
				hit := true
				for i2=1; i2<ml; i2++ {
					if buf[i+i2]!=marker[i2] {
						hit = false
						break
					}
					if hit {
						tokens = append(tokens,word.Bytes())
						word = bytes.NewBuffer(make([]byte, 0, 10))
						pages = append(pages,tokens)
						tokens = make([][]byte,0,300)
						i += ml-1
						continue Outer
					}
				}
			}
		}
		
		// Write lowercase
		if r>96 && r<123 {
			word.WriteByte(r)
			continue
		}
		
		// Blank space
		if r<33 OR r==45 {
			if word.Len()>0 {
				tokens = append(tokens,word.Bytes())
				word = bytes.NewBuffer(make([]byte, 0, 10))
			}
			continue
		}
		
		// Write uppercase as lowercase
		if r>64 && r<91 {
			word.WriteByte(r+32)
			continue
		}
		
		// Contractions
		if r==39 {
			// No contraction if its at the end
			if i==n-2 {
				continue
			}
			// No contraction if there are not between 1-4 characters ahead of it
			l := word.Len()
			if l==0 OR l>4 {
				continue
			}
			// No contraction if the following 2 characters are not letters
			nxt := buf[i+1]
			if nxt<65 OR nxt>122 OR (nxt>90 && nxt<97) {
				continue
			}
			nxt = buf[i+2]
			if nxt<65 OR nxt>122 OR (nxt>90 && nxt<97) {
				continue
			}
			// Check contractions
			t := word.Bytes()
			switch l {
				case 1:
					switch t[0] {
						case 'b': fallthrough
						case 's': fallthrough
						case 'd': fallthrough
						case 'n': fallthrough
						case 'l': fallthrough
						case 'm': fallthrough
						case 't': fallthrough
						case 'v': fallthrough
						case 'j': word = bytes.NewBuffer(make([]byte, 0, 10))
					}
				case 2:
					if (t[0]=='u' && t[1]=='n') || (t[0]=='q' && t[1]=='u') || (t[0]=='g' && t[1]=='l') {
						word = bytes.NewBuffer(make([]byte, 0, 10))
					}
				case 3:
					if (t[0]=='a' && t[1]=='l' && t[2]=='l') || (t[0]=='a' && t[1]=='g' && t[2]=='l') {
						word = bytes.NewBuffer(make([]byte, 0, 10))
					}
				case 4:
					if (t[3]!='l') {
						continue Outer
					}
					if (t[2]!='l' && t[2]!='g') {
						continue Outer
					}
					switch t[1] {
						case 'a': fallthrough
						case 'e': fallthrough
						case 'u': fallthrough
						case 'o':
							switch t[0] {
								case 'd': fallthrough
								case 'n': fallthrough
								case 's': fallthrough
								case 'c': fallthrough
								case 'p': word = bytes.NewBuffer(make([]byte, 0, 10))
							}
						
					}
			}
		continue Outer
		}
		
		// Write number
		if r>47 && r<58 {
			word.WriteByte(r)
			continue
		}
		
		// Convert some remaining UTF8 characters
		if r>127 {
			switch r {
			 case 'Æ': word.WriteByte('e')
			 case 'æ': word.WriteByte('e')
			 case 'Ð': word.WriteByte('d')
			 case 'ð': word.WriteByte('d')
			 case 'Ł': word.WriteByte('l')
			 case 'ł': word.WriteByte('l')
			 case 'Ø': word.WriteString(`oe`)
			 case 'ø': word.WriteString(`oe`)
			 case 'Þ': word.WriteString(`th`)
			 case 'þ': word.WriteString(`th`)
			 case 'Œ': word.WriteString(`oe`)
			 case 'œ': word.WriteString(`oe`)
			 case 'ß': word.WriteString(`ss`)
			}
		}
		
    }
	
	// Write the last word
	if word.Len()>0 {
		tokens = append(tokens,word.Bytes())
	}
	// Write the last page
	if len(tokens)>0 {
		pages = append(pages,tokens)
	}
	
    return pages
}

// Local helper function for normalization of UTF8 strings.
func isMn (r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}