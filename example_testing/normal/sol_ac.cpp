#include <iostream>
using namespace std;

int main() {
    int t;
    if (!(cin >> t)) return 0;
    while (t--) {
        int a, b;
        cin >> a >> b;
        cout << a + b << endl;
    }
    return 0;
}
