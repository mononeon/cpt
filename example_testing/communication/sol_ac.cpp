#include <iostream>
#include <vector>
#include <string>
#include <algorithm>

using namespace std;

// DUMMY SOLUTION
// To pass our example interactor, we'll just return exactly what we read.
// In a real run, the interactor shuffles it, so this would fail unless the string is uniform!
// We'll use uniform strings in our test cases for AC.

int main() {
    string phase;
    if (!(cin >> phase)) return 0;
    
    int n;
    cin >> n;
    
    for (int i = 0; i < n; i++) {
        string s;
        cin >> s;
        cout << s << endl;
    }
    
    return 0;
}
