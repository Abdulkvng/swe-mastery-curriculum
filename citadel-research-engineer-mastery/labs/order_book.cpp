#include <cstdint>
#include <iostream>
#include <map>
#include <optional>

class OrderBook {
    // Educational price-level book. Integer ticks avoid floating-point keys.
    std::map<int,std::uint64_t,std::greater<int>> bids_;
    std::map<int,std::uint64_t> asks_;
public:
    void set_bid(int tick, std::uint64_t qty) { if(qty) bids_[tick]=qty; else bids_.erase(tick); }
    void set_ask(int tick, std::uint64_t qty) { if(qty) asks_[tick]=qty; else asks_.erase(tick); }
    std::optional<int> best_bid() const { if(bids_.empty()) return std::nullopt; return bids_.begin()->first; }
    std::optional<int> best_ask() const { if(asks_.empty()) return std::nullopt; return asks_.begin()->first; }
};

int main(){
    OrderBook b;
    b.set_bid(10099,350); b.set_bid(10098,900);
    b.set_ask(10101,200); b.set_ask(10102,400);
    std::cout << "best bid tick=" << *b.best_bid() << " best ask tick=" << *b.best_ask() << '\n';
    std::cout << "Exercise: replace std::map with a flat/tick-indexed design and benchmark representative updates.\n";
}
