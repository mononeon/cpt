#include<iostream>
#include<set>
using namespace std;
int t,n;
set<long long> s;
int makeI(long long x)
{
	cout<<"I "<<x<<endl;
	int ans=0;
	cin>>ans;
	return ans;
}
int makeQ(long long y)
{
	cout<<"Q "<<y<<endl;
	int ans=0;
	cin>>ans;
	return ans;
}
int query(long long x)
{
	int res=makeQ(x);
	for(long long w:s)
	if(w>=x)
	 res--;
	return res;
}
long long guess(int n)
{
	long long pre=0;
	for(int i=n-1;i>=0;i--)
	if(query(pre|(1ll<<i))>0)
	{
	 	pre|=(1ll<<i);
	}
	return pre;
}
int main()
{
	cin>>t;
	while(t--)
	{
		cin>>n;
		s.clear();
		cout<<0<<endl;
		s.insert(0);
		int res1=makeI(0);
		if(res1==1)
		{
			makeI((1ll<<n)-1);
			long long w=guess(n);
			cout<<"A "<<1<<" "<<w<<endl;
			continue;
		}
		long long c=guess(n);
		s.insert(c);
		if(c==(1ll<<n)-1)
		{
			int res2=makeI(1);
			if(res2!=res1)
			{
				cout<<"A "<<3<<" "<<c<<endl;
			}
			else
			{
				cout<<"A "<<2<<" "<<c<<endl;
			}
			continue;
		}
		int res2=makeI((1ll<<n)-1);
		if(query((1ll<<n)-1)>0)
		{
			cout<<"A "<<2<<" "<<c<<endl;
		}
		else
		{
			cout<<"A "<<3<<" "<<c<<endl;
		}
	}
}
