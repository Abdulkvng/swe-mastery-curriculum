#include <algorithm>
#include <chrono>
#include <cstdint>
#include <iostream>
#include <numeric>
#include <random>
#include <vector>

using Clock = std::chrono::steady_clock;

static volatile std::uint64_t sink = 0;

template <class F>
double ns_per_access(F&& fn, std::size_t accesses) {
    auto t0 = Clock::now();
    sink ^= fn();
    auto t1 = Clock::now();
    return std::chrono::duration<double, std::nano>(t1 - t0).count() / accesses;
}

int main() {
    constexpr std::size_t n = 8 * 1024 * 1024;
    std::vector<std::uint64_t> data(n);
    std::iota(data.begin(), data.end(), 1);

    std::vector<std::size_t> idx(n);
    std::iota(idx.begin(), idx.end(), 0);
    std::mt19937_64 rng(42);
    std::shuffle(idx.begin(), idx.end(), rng);

    auto sequential = [&] {
        std::uint64_t sum = 0;
        for (auto x : data) sum += x;
        return sum;
    };
    auto random = [&] {
        std::uint64_t sum = 0;
        for (auto i : idx) sum += data[i];
        return sum;
    };

    std::cout << "sequential ns/access: " << ns_per_access(sequential, n) << '\n';
    std::cout << "random     ns/access: " << ns_per_access(random, n) << '\n';
    std::cout << "sink=" << sink << '\n';
}
