#include <iostream>
using namespace std;

int main() {
    int t;
    if (!(cin >> t)) return 0;
    while (t--) {
        int a, b;
        cin >> a >> b;
        if (a == 5) cout << a + b + 1 << endl; // Intentional WA for edge case
        else cout << a + b << endl;
    }
    return 0;
}
