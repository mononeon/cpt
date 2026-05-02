#include "testlib.h"
#include <iostream>

using namespace std;

int main(int argc, char* argv[]) {
    registerGen(argc, argv, 1);
    
    int t = rnd.next(1, 5);
    cout << t << endl;
    
    for (int i = 0; i < t; i++) {
        // We force '5' to appear sometimes to trigger our intentional WA/TLE/RE
        int a = rnd.next(1, 10);
        int b = rnd.next(1, 10);
        cout << a << " " << b << endl;
    }
    
    return 0;
}
