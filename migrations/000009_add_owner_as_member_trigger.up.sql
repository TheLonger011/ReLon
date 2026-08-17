CREATE FUNCTION add_owner_as_member()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO community_members (community_id,user_id)
    VALUES (NEW.id,NEW.owner_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_add_owner_as_member
AFTER INSERT ON communities
FOR EACH ROW
EXECUTE FUNCTION add_owner_as_member();