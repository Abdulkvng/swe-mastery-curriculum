#include <array>
#include <atomic>
#include <cassert>
#include <cstddef>
#include <iostream>
#include <thread>

template<class T, std::size_t N>
class SPSCQueue {
    static_assert(N >= 2);
    std::array<T,N> data_{};
    alignas(64) std::atomic<std::size_t> head_{0}; // producer writes
    alignas(64) std::atomic<std::size_t> tail_{0}; // consumer writes
public:
    bool push(const T& v) {
        const auto h=head_.load(std::memory_order_relaxed);
        const auto next=(h+1)%N;
        if(next==tail_.load(std::memory_order_acquire)) return false;
        data_[h]=v;
        head_.store(next,std::memory_order_release);
        return true;
    }
    bool pop(T& out) {
        const auto t=tail_.load(std::memory_order_relaxed);
        if(t==head_.load(std::memory_order_acquire)) return false;
        out=data_[t];
        tail_.store((t+1)%N,std::memory_order_release);
        return true;
    }
};

int main(){
    constexpr int count=1'000'000;
    SPSCQueue<int,1024> q;
    long long sum=0;
    std::thread producer([&]{ for(int i=1;i<=count;){ if(q.push(i)) ++i; } });
    std::thread consumer([&]{ int v; for(int i=0;i<count;){ if(q.pop(v)){ sum+=v; ++i; } } });
    producer.join(); consumer.join();
    assert(sum == 1LL*count*(count+1)/2);
    std::cout << "ok sum=" << sum << '\n';
}
