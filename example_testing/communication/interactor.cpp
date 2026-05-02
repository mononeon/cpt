#include "testlib.h"
#include <iostream>
#include <vector>

using namespace std;

int main(int argc, char* argv[]) {
    registerInteraction(argc, argv);
    
    int n = inf.readInt();
    vector<string> combs(n);
    for (int i = 0; i < n; i++) combs[i] = inf.readToken();
    
    // --- PHASE 1: ENCODE ---
    cout << "encode" << endl;
    cout << n << endl;
    for (int i = 0; i < n; i++) cout << combs[i] << endl;
    
    vector<string> encodings(n);
    for (int i = 0; i < n; i++) {
        encodings[i] = ouf.readToken();
        if (encodings[i].length() > 67) {
            quitf(_wa, "Encoding length exceeded 67 digits");
        }
    }
    
    // --- PHASE 2: DECODE ---
    cout << "decode" << endl;
    cout << n << endl;
    for (int i = 0; i < n; i++) {
        string shuffled = encodings[i];
        shuffle(shuffled.begin(), shuffled.end()); // testlib shuffle
        cout << shuffled << endl;
    }
    
    for (int i = 0; i < n; i++) {
        string dec = ouf.readToken();
        if (dec != combs[i]) {
            quitf(_wa, "Wrong decode: expected %s, found %s", combs[i].c_str(), dec.c_str());
        }
    }
    
    quitf(_ok, "Successfully encoded and decoded %d combinations", n);
}
