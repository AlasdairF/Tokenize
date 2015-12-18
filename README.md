##Tokenize

This Tokenize package contains three functions that are extremely fast and efficient at tokenizing text. No regular expressions are used. The whole thing requires only two loops of the data, the first for UTF8 normalization and accent removal, the second for everything else.

##Warning
The same underlying array is used for each token, this means you **must** copy the slice of bytes sent to the `wordfn` function if you intend to save the slices. Please see my [Unleak](http://github.com/AlasdairF/Unleak) package for an easy one-liner implementation of this. If you are counting the token occurances with my [BinSearch](http://github.com/AlasdairF/BinSearch) package, with the native `map` implementation, or you are converting the slice of bytes to a string then it is not necessary to copy the slice since these implementations make their own copies.

##Features

 * UTF8 normalization
 * Accent removal, e.g. á -> a
 * Special characters converted to their common form, e.g. æ -> e
 * Lowercasing
 * Hypenated words split
 * Contractions removed, e.g. l'histoire -> histoire (but they're -> theyre)
 * All UTF8 scripts are supported.

For example:

    Et l'Histore de l'amitè.
	
Becomes

    et
    histore
    de
    amite

##Installation

    go get github.com/AlasdairF/Tokenize

##Parameters

The optional parameters are:

    lowercase				converts all letters to lowercase

    stripAccents			removes accents, e.g. á -> a

    stripContractions		removes contractions, e.g. l'histore -> histore

    stripNumbers			removes all numbers

    stripForeign			leaves only a-zA-Z0-9 (after accent removal)

##Recommended

Recommended settings for tokenization of English are:

    lowercase, stripAccents, stripForeign

Recommended settings for tokenization of continental European languages are:

    lowercase, stripAccents, stripContractions, stripForeign

Recommended settings for tokenization of international scripts are:

    lowercase, stripContractions

All non-letters and non-numbers, such as punctuation and whitespace are always stripped.

##AllInOne

The first parameter is the []byte data to process, the second is the function for what to with each token. Then the options.

For example, if you want to put all words into a slice then you would use:

    tokens := make([][]byte, 0, 100)
    
    wordfn := func(word []byte) {
    	tokens = append(tokens, unleak.Bytes(word)) // using my Unleak package to make a copy of the slice
    }
    
    lowercase, stripAccents, stripContractions, stripNumbers, stripForeign := true, true, true, false, true
    tokenize.AllInOne(data, wordfn, lowercase, stripAccents, stripContractions, stripNumbers, stripForeign)

## WithProvidedBuffer

Exactly the same as AllInOne but accepts the reuse of the `custom.Buffer`. This is much faster if you are repeatedly using this package on small chunks of data.

    import "github.com/AlasdairF/Custom"
    buf := custom.NewBuffer(32)
    tokenize.WithProvidedBuffer(buf, data, wordfn, lowercase, stripAccents, stripContractions, stripNumbers, stripForeign)

##Paginate

Paginate is the same as AllInOne but it also recognizes custom page breaks. Four parameters are required. The first is the []byte data to process, the second is the page break marker as []byte, the third is the function for what to do with each token, the fourth is the function for what to do whenever a page break marker is reached. Please note that the page break marker itself should contain only single-byte characters (ASCII), I usually use `{#}` as the page break marker.

For example:

    pages := make([][][]byte, 0, 10)
    tokens := make([][]byte, 0, 100)
	
	wordfn := func(word []byte) {
    	tokens = append(tokens, unleak.Bytes(word))
    }
	pagefn := func() {
		pages = append(pages, tokens)
		tokens = make([][]byte, 0, 100)
    }
    
    lowercase, stripAccents, stripContractions, stripNumbers, stripForeign := true, true, true, false, true
    tokenize.Paginate(data, []byte("[newpage]"), wordfn, pagefn, lowercase, stripAccents, stripContractions, stripNumbers, stripForeign)
