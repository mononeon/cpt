#include <iostream>
#include <vector>
#include <string>

using namespace std;

int main() {
    string phase;
    if (!(cin >> phase)) return 0;
    
    int n;
    cin >> n;
    
    for (int i = 0; i < n; i++) {
        string s;
        cin >> s;
        cout << "12345" << endl; // Wrong output
    }
    
    return 0;
}
