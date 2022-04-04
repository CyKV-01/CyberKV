use sonyflake::Sonyflake;

pub fn next_id() -> u64 {
    let mut sonyflake = Sonyflake::new().unwrap();
    sonyflake.next_id().unwrap()
}
