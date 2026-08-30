#include <algorithm>
#include <chrono>
#include <cstdint>
#include <iostream>
#include <random>
#include <vector>

using Clock = std::chrono::steady_clock;
static volatile std::uint64_t sink = 0;

std::uint64_t count_branch(const std::vector<int>& v) {
    std::uint64_t c = 0;
    for (int x : v) if (x >= 128) ++c;
    return c;
}

int main() {
    constexpr std::size_t n = 20'000'000;
    std::vector<int> random_data(n);
    std::mt19937 rng(42);
    std::uniform_int_distribution<int> dist(0, 255);
    for (auto& x : random_data) x = dist(rng);

    auto sorted_data = random_data;
    std::sort(sorted_data.begin(), sorted_data.end());

    for (auto* p : {&sorted_data, &random_data}) {
        auto t0 = Clock::now();
        sink ^= count_branch(*p);
        auto t1 = Clock::now();
        std::cout << std::chrono::duration<double, std::milli>(t1-t0).count() << " ms\n";
    }
    std::cout << "Run with Linux perf stat -e branches,branch-misses ./branch_benchmark\n";
}
