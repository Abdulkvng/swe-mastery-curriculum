#include <atomic>
#include <chrono>
#include <cstdint>
#include <iostream>
#include <thread>

constexpr std::uint64_t ITERS = 50'000'000;

struct Packed { std::atomic<std::uint64_t> a{0}, b{0}; };
struct alignas(64) PaddedCounter { std::atomic<std::uint64_t> x{0}; };
struct Padded { PaddedCounter a, b; };

template<class T>
double run(T& s) {
    auto t0 = std::chrono::steady_clock::now();
    std::thread t1([&]{ for (std::uint64_t i=0;i<ITERS;++i) s.a.x.fetch_add(1,std::memory_order_relaxed); });
    std::thread t2([&]{ for (std::uint64_t i=0;i<ITERS;++i) s.b.x.fetch_add(1,std::memory_order_relaxed); });
    t1.join(); t2.join();
    return std::chrono::duration<double>(std::chrono::steady_clock::now()-t0).count();
}

int main() {
    // The packed version is written separately because fields differ in shape.
    Packed p;
    auto t0=std::chrono::steady_clock::now();
    std::thread p1([&]{for(std::uint64_t i=0;i<ITERS;++i)p.a.fetch_add(1,std::memory_order_relaxed);});
    std::thread p2([&]{for(std::uint64_t i=0;i<ITERS;++i)p.b.fetch_add(1,std::memory_order_relaxed);});
    p1.join();p2.join();
    double packed=std::chrono::duration<double>(std::chrono::steady_clock::now()-t0).count();

    Padded q;
    std::cout << "packed seconds: " << packed << '\n';
    std::cout << "padded seconds: " << run(q) << '\n';
}
