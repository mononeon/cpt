#include <iostream>
using namespace std;

int main() {
    int t;
    if (!(cin >> t)) return 0;
    while (t--) {
        int a, b;
        cin >> a >> b;
        if (a == 5) {
            while (true) {} // Intentional TLE
        }
        cout << a + b << endl;
    }
    return 0;
}
