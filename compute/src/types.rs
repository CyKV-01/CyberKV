pub type SlotID = u16;
pub type TimeStamp = u64;

pub const SLOT_NUM: u16 = 32;

#[derive(Clone, Debug)]
pub struct Value {
    pub timestamp: TimeStamp,
    pub value: String,
}

impl Value {
    pub fn new(ts: TimeStamp, value: String) -> Self {
        return Self {
            timestamp: ts,
            value: value,
        };
    }
}
