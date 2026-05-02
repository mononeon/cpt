package templates

import (
	"fmt"
	"os"
)

const genTemplate = `// Generator Template
#include "testlib.h"
#include <iostream>

using namespace std;

int main(int argc, char* argv[]) {
    registerGen(argc, argv, 1);
    
    // Example: Generate N between 1 and 10
    int n = rnd.next(1, 10);
    cout << n << endl;
    
    for (int i = 0; i < n; i++) {
        cout << rnd.next(1, 100) << (i == n - 1 ? "" : " ");
    }
    cout << endl;
    
    return 0;
}
`

const checkTemplate = `// Checker Template
#include "testlib.h"
#include <iostream>

using namespace std;

int main(int argc, char* argv[]) {
    // argv[1]: Input, argv[2]: Output, argv[3]: Answer
    registerTestlibCmd(argc, argv);
    
    int expected = ans.readInt();
    int actual = ouf.readInt();
    
    if (expected == actual) {
        quitf(_ok, "Answers match: %d", actual);
    } else {
        quitf(_wa, "Expected %d, found %d", expected, actual);
    }
}
`

const valTemplate = `// Validator Template
#include "testlib.h"
#include <iostream>

using namespace std;

int main(int argc, char* argv[]) {
    registerValidation(argc, argv);
    
    // Example: N between 1 and 10
    int n = inf.readInt(1, 10, "n");
    inf.readEoln();
    
    for (int i = 0; i < n; i++) {
        inf.readInt(1, 100, "a_i");
        if (i < n - 1) inf.readSpace();
    }
    inf.readEoln();
    inf.readEof();
    
    return 0; // Exits with 0 if valid
}
`

const interactTemplate = `// Interactor Template
#include "testlib.h"
#include <iostream>

using namespace std;

int main(int argc, char* argv[]) {
    // argv[1]: Input, argv[2]: Output, argv[3]: Answer
    registerInteraction(argc, argv);
    
    // Read secret number from input file
    int secret = inf.readInt();
    
    int queries = 0;
    while (true) {
        string type = ouf.readString();
        if (type == "!") {
            int guess = ouf.readInt();
            if (guess == secret) {
                quitf(_ok, "Guessed correctly in %d queries", queries);
            } else {
                quitf(_wa, "Wrong guess. Expected %d, got %d", secret, guess);
            }
        } else if (type == "?") {
            queries++;
            if (queries > 20) quitf(_wa, "Too many queries");
            
            int q = ouf.readInt();
            if (q < secret) cout << "<" << endl;
            else if (q > secret) cout << ">" << endl;
            else cout << "=" << endl;
        } else {
            quitf(_pe, "Unknown query type");
        }
    }
}
`

func GenerateTemplate(t string) {
	var content, filename string

	switch t {
	case "gen":
		content = genTemplate
		filename = "gen.cpp"
	case "check":
		content = checkTemplate
		filename = "checker.cpp"
	case "val":
		content = valTemplate
		filename = "validator.cpp"
	case "interact":
		content = interactTemplate
		filename = "interactor.cpp"
	default:
		fmt.Println("\033[1;31mUnknown template type. Use: gen, check, val, interact\033[0m")
		return
	}

	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("\033[33m%s already exists. Skipping...\033[0m\n", filename)
		return
	}

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to create template: %v\033[0m\n", err)
		return
	}

	fmt.Printf("\033[32mSuccessfully created %s\033[0m\n", filename)
}
