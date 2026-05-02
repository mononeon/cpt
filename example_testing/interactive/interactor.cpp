#include "testlib.h"
#include <iostream>
#include <set>

using namespace std;

int main(int argc, char* argv[]) {
    registerInteraction(argc, argv);
    
    int t = inf.readInt();
    cout << t << endl;
    
    while (t--) {
        int n = inf.readInt();
        int k = inf.readInt();
        long long c = inf.readLong();
        
        cout << n << endl;
        
        long long a = ouf.readLong();
        set<long long> S;
        S.insert(a);
        
        int queries = 0;
        bool answered = false;
        
        while (queries <= n + 3) {
            string type = ouf.readToken();
            if (type == "I") {
                queries++;
                long long x = ouf.readLong();
                long long fx = 0;
                if (k == 1) fx = x & c;
                else if (k == 2) fx = x | c;
                else if (k == 3) fx = x ^ c;
                
                S.insert(fx);
                cout << S.size() << endl;
            } else if (type == "Q") {
                queries++;
                long long y = ouf.readLong();
                int cnt = 0;
                for (long long z : S) {
                    if (z >= y) cnt++;
                }
                cout << cnt << endl;
            } else if (type == "A") {
                int ans_k = ouf.readInt();
                long long ans_c = ouf.readLong();
                if (ans_k == k && ans_c == c) {
                    answered = true;
                    break;
                } else {
                    quitf(_wa, "Wrong Answer: Expected %d %lld, got %d %lld", k, c, ans_k, ans_c);
                }
            } else {
                quitf(_pe, "Unknown query type: %s", type.c_str());
            }
        }
        
        if (!answered) {
            quitf(_wa, "Did not answer within %d queries", n + 3);
        }
    }
    
    quitf(_ok, "All %d test cases answered correctly", t);
}
