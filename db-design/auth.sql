DROP TABLE IF EXISTS auth CASCADE;
CREATE TABLE auth
(
    auth_id int8 NOT NULL GENERATED ALWAYS AS IDENTITY (INCREMENT 1 MINVALUE 1 NO MAXVALUE START 1),
    lang VARCHAR(3) DEFAULT 'tr' NOT NULL,
    src INT2 NOT NULL DEFAULT 1 ,
    manager_id INT8 NOT NULL DEFAULT 0 ,
    user_type INT2 NOT NULL DEFAULT 2,
    user_gender INT2 NOT NULL DEFAULT 0,
    user_title VARCHAR(255) NOT NULL,
    user_name VARCHAR(65) NOT NULL,
    user_pass VARCHAR(65) NOT NULL,
    staff_id INT8 DEFAULT 0 NOT NULL CHECK (staff_id >= 0),
    unique_id uuid NOT NULL DEFAULT gen_random_uuid(),
    status INT2 DEFAULT 4 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    expired_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    PRIMARY KEY (auth_id),
    UNIQUE (user_name)

);
CREATE INDEX ON auth USING btree (lang);
CREATE INDEX ON auth USING btree (user_title);
CREATE INDEX ON auth USING btree (user_name);
CREATE INDEX ON auth USING btree (staff_id);
COMMENT ON COLUMN auth.src is '1) Web, 2) Mobile, 3) Other';
COMMENT ON COLUMN auth.user_type is '1) Hire Expert, 2) Expert';
COMMENT ON COLUMN auth.user_name is 'User Mail';
COMMENT ON COLUMN auth.status is '1) Active, 2) Passive, 3) On Hold, 4) Waiting Verification by Authorize, 5) Suspended';


DROP TABLE IF EXISTS notifications CASCADE;
CREATE TABLE notifications
(
    id int8 NOT NULL GENERATED ALWAYS AS IDENTITY (INCREMENT 1 MINVALUE 1 NO MAXVALUE START 1),
    auth_id INT8 NOT NULL UNIQUE,
    is_act_mail BOOLEAN DEFAULT TRUE NOT NULL,
    is_act_phone BOOLEAN DEFAULT TRUE NOT NULL,
    is_act_push BOOLEAN DEFAULT TRUE NOT NULL,
    is_act_2fa BOOLEAN DEFAULT false NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX ON notifications USING btree (auth_id);


DROP TABLE IF EXISTS auth_avatar CASCADE;
CREATE TABLE auth_avatar
(
    auth_id INT8 NOT NULL UNIQUE,
    image_url VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (auth_id)
);
CREATE INDEX ON auth_avatar USING btree (auth_id);


DROP TABLE IF EXISTS auth_address CASCADE;
CREATE TABLE auth_address
(
    id int8 NOT NULL GENERATED ALWAYS AS IDENTITY (INCREMENT 1 MINVALUE 1 NO MAXVALUE START 1),
    auth_id INT8 NOT NULL UNIQUE,
    description VARCHAR(75) DEFAULT 0,
    definition VARCHAR(255) DEFAULT 0,
    country_id INT8 DEFAULT 0,
    country_name VARCHAR(50) DEFAULT NULL,
    city_id INT8 DEFAULT 0,
    city_name VARCHAR(50) DEFAULT NULL,
    town_id INT8 DEFAULT 0,
    town_name VARCHAR(50) DEFAULT NULL,
    district_id INT8 DEFAULT 0,
    district_name VARCHAR(50) DEFAULT NULL,
    quarter_id INT8 DEFAULT 0,
    quarter_name VARCHAR(50) DEFAULT NULL,
    zip_code VARCHAR(10) DEFAULT NULL,
    is_dafault BOOLEAN NOT NULL DEFAULT false,
    status INT2 NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX ON auth_address USING btree (auth_id);
CREATE INDEX ON auth_address USING btree (description);
CREATE INDEX ON auth_address USING btree (country_id);
CREATE INDEX ON auth_address USING btree (city_id);
CREATE INDEX ON auth_address USING btree (town_id);
CREATE INDEX ON auth_address USING btree (district_id);
CREATE INDEX ON auth_address USING btree (quarter_id);
CREATE INDEX ON auth_address USING btree (quarter_id);


DROP TABLE IF EXISTS auth_contact CASCADE;
CREATE TABLE auth_contact
(
    id int8 NOT NULL GENERATED ALWAYS AS IDENTITY (INCREMENT 1 MINVALUE 1 NO MAXVALUE START 1),
    auth_id INT8 NOT NULL UNIQUE,
    type_of INT2 DEFAULT 1,
    description VARCHAR(75) DEFAULT 0,
    definition VARCHAR(255) DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    status INT2 NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX ON auth_contact USING btree (auth_id);
CREATE INDEX ON auth_contact USING btree (type_of);
CREATE INDEX ON auth_contact USING btree (type_of);
COMMENT ON COLUMN auth_contact.type_of is '1) Mobile, 2) Work Phone, 3) E-Mail, 4) Web Site, 5) Location, 6) Fax, 7) Social Media';


DROP TABLE IF EXISTS user_tokens CASCADE;
CREATE TABLE user_tokens
(
    id INT8 NOT NULL GENERATED ALWAYS AS IDENTITY (INCREMENT 1 MINVALUE 1 NO MAXVALUE START 1),
    auth_id INT8 NOT NULL UNIQUE CHECK (auth_id > 0),
    token VARCHAR[] DEFAULT NULL,
    topics VARCHAR[] DEFAULT NULL,
    device_os INT2 DEFAULT NULL,
    device_type INT2 NOT NULL DEFAULT 0,
    browser_type INT2 NOT NULL DEFAULT 0,
    status INT2 NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX ON user_tokens USING btree (auth_id);
CREATE INDEX ON user_tokens USING btree (token);
CREATE INDEX ON user_tokens USING btree (device_os);
COMMENT ON COLUMN user_tokens.token IS 'Generated Token by FCM api, its unique for every device with browser';
COMMENT ON COLUMN user_tokens.topics IS 'Subscribed topics';
COMMENT ON COLUMN user_tokens.device_os IS '1) Android, 2) iOS, 3) HarmonyOs, 4) Win, 5) Mac, 6) Linux';
COMMENT ON COLUMN user_tokens.device_type IS '1) Mobile, 2) Desktop';
COMMENT ON COLUMN user_tokens.browser_type IS '1) Chrome, 2) Firefox, 3) Edge, 4) Opera, 5) Brave, 6) Safari';
