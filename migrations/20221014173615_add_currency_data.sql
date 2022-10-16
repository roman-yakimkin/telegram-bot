-- +goose Up
-- +goose StatementBegin
INSERT INTO currency VALUES ('USD', '$%.2f'),('CNY', '%.2f元'),('EUR', '%.2f€'),('RUB', '%.2f₽');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
