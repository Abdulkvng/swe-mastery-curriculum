// LRU cache in safe Rust.
//
// Why this exercise teaches Rust:
//   - Doubly-linked lists in Rust are notoriously tricky because each node
//     has TWO owners (prev pointer + next pointer + the map). Plain ownership
//     doesn't allow that.
//   - Solution 1 (this file): use Rc<RefCell<Node>> — runtime-checked shared
//     mutability. Easier to reason about, slight runtime cost.
//   - Solution 2 (alternative not shown): unsafe raw pointers. Faster, but
//     you become responsible for memory safety yourself. Most production
//     Rust crates use this.
//
// Rc = Reference Counted: lets multiple variables co-own a value.
//      Single-threaded only. (Use Arc for multi-threaded.)
// RefCell = "interior mutability": borrow-check at runtime, not compile-time.
//      You can mutate through an immutable reference.

use std::cell::RefCell;
use std::collections::HashMap;
use std::hash::Hash;
use std::rc::Rc;

type Link<K, V> = Option<Rc<RefCell<Node<K, V>>>>;

struct Node<K, V> {
    key: K,
    value: V,
    prev: Link<K, V>,
    next: Link<K, V>,
}

pub struct LruCache<K, V>
where
    K: Hash + Eq + Clone,
{
    capacity: usize,
    map: HashMap<K, Rc<RefCell<Node<K, V>>>>,
    // Sentinel head and tail. head.next = MRU, tail.prev = LRU.
    head: Rc<RefCell<Node<K, V>>>,
    tail: Rc<RefCell<Node<K, V>>>,
}

impl<K, V> LruCache<K, V>
where
    K: Hash + Eq + Clone,
    V: Clone,
{
    pub fn new(capacity: usize) -> Self {
        assert!(capacity > 0, "capacity must be > 0");

        // Sentinels need *some* K and V. We use a hack: leave them
        // logically uninitialized via MaybeUninit-equivalent... Actually,
        // for this teaching example, we'll require K: Default + V: Default
        // — see new_with_default below. To avoid bounds bloat, we go a
        // different route: store Option<K>/Option<V> in nodes. But that
        // complicates the rest. Cleanest: split internal node from
        // public node, with Option fields for sentinels.
        //
        // For brevity we expose only `new_with_defaults` which requires Default.
        unimplemented!("use LruCache::with_defaults")
    }
}

// To keep the example compileable, we provide a version that works when
// K and V both implement Default — a totally reasonable trade-off for
// teaching the structure clearly. Real production crates handle this with
// raw pointers or Option<...> fields.

impl<K, V> LruCache<K, V>
where
    K: Hash + Eq + Clone + Default,
    V: Clone + Default,
{
    pub fn with_defaults(capacity: usize) -> Self {
        assert!(capacity > 0, "capacity must be > 0");
        let head = Rc::new(RefCell::new(Node {
            key: K::default(),
            value: V::default(),
            prev: None,
            next: None,
        }));
        let tail = Rc::new(RefCell::new(Node {
            key: K::default(),
            value: V::default(),
            prev: None,
            next: None,
        }));
        head.borrow_mut().next = Some(tail.clone());
        tail.borrow_mut().prev = Some(head.clone());
        LruCache {
            capacity,
            map: HashMap::with_capacity(capacity),
            head,
            tail,
        }
    }

    pub fn len(&self) -> usize {
        self.map.len()
    }

    pub fn is_empty(&self) -> bool {
        self.map.is_empty()
    }

    pub fn get(&mut self, key: &K) -> Option<V> {
        if let Some(node) = self.map.get(key).cloned() {
            self.move_to_front(&node);
            Some(node.borrow().value.clone())
        } else {
            None
        }
    }

    pub fn put(&mut self, key: K, value: V) {
        if let Some(node) = self.map.get(&key).cloned() {
            node.borrow_mut().value = value;
            self.move_to_front(&node);
            return;
        }

        let new_node = Rc::new(RefCell::new(Node {
            key: key.clone(),
            value,
            prev: None,
            next: None,
        }));
        self.add_to_front(&new_node);
        self.map.insert(key, new_node);

        if self.map.len() > self.capacity {
            // Evict node before tail.
            let lru = self.tail.borrow().prev.clone().unwrap();
            self.unlink(&lru);
            let lru_key = lru.borrow().key.clone();
            self.map.remove(&lru_key);
        }
    }

    pub fn remove(&mut self, key: &K) -> bool {
        if let Some(node) = self.map.remove(key) {
            self.unlink(&node);
            true
        } else {
            false
        }
    }

    // Insert `n` right after head.
    fn add_to_front(&self, n: &Rc<RefCell<Node<K, V>>>) {
        let head_next = self.head.borrow().next.clone().unwrap();

        n.borrow_mut().prev = Some(self.head.clone());
        n.borrow_mut().next = Some(head_next.clone());

        self.head.borrow_mut().next = Some(n.clone());
        head_next.borrow_mut().prev = Some(n.clone());
    }

    // Remove `n` from its position. Doesn't touch the map.
    fn unlink(&self, n: &Rc<RefCell<Node<K, V>>>) {
        let prev = n.borrow().prev.clone();
        let next = n.borrow().next.clone();
        if let (Some(p), Some(nx)) = (prev, next) {
            p.borrow_mut().next = Some(nx.clone());
            nx.borrow_mut().prev = Some(p);
        }
        n.borrow_mut().prev = None;
        n.borrow_mut().next = None;
    }

    fn move_to_front(&self, n: &Rc<RefCell<Node<K, V>>>) {
        self.unlink(n);
        self.add_to_front(n);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn basic() {
        let mut c: LruCache<String, i32> = LruCache::with_defaults(2);
        c.put("a".into(), 1);
        c.put("b".into(), 2);
        assert_eq!(c.get(&"a".into()), Some(1));
        c.put("c".into(), 3); // evicts b
        assert_eq!(c.get(&"b".into()), None);
        assert_eq!(c.get(&"c".into()), Some(3));
    }

    #[test]
    fn update_existing() {
        let mut c: LruCache<String, i32> = LruCache::with_defaults(2);
        c.put("a".into(), 1);
        c.put("b".into(), 2);
        c.put("a".into(), 99); // bumps a to front
        c.put("c".into(), 3);  // evicts b, not a
        assert_eq!(c.get(&"a".into()), Some(99));
        assert_eq!(c.get(&"b".into()), None);
    }

    #[test]
    fn remove() {
        let mut c: LruCache<String, i32> = LruCache::with_defaults(3);
        c.put("a".into(), 1);
        c.put("b".into(), 2);
        c.put("c".into(), 3);
        assert!(c.remove(&"b".into()));
        assert_eq!(c.len(), 2);
        assert_eq!(c.get(&"b".into()), None);
    }

    #[test]
    fn evicts_lru() {
        let mut c: LruCache<i32, i32> = LruCache::with_defaults(3);
        c.put(1, 1);
        c.put(2, 2);
        c.put(3, 3);
        c.get(&1);   // 1 is now MRU; LRU = 2
        c.put(4, 4); // evict 2
        assert_eq!(c.get(&2), None);
        assert_eq!(c.get(&1), Some(1));
        assert_eq!(c.get(&3), Some(3));
        assert_eq!(c.get(&4), Some(4));
    }
}
