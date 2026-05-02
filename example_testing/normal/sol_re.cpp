#include <iostream>
using namespace std;

int main() {
    int t;
    if (!(cin >> t)) return 0;
    while (t--) {
        int a, b;
        cin >> a >> b;
        if (a == 5) {
            int x = 0;
            cout << a / x << endl; // Intentional RE (Division by zero)
        }
        cout << a + b << endl;
    }
    return 0;
}
