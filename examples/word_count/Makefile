BUCKET = ${AWS_TEST_BUCKET}
BIN_DIR = ./bin
PROG_NAME = word_count
HOST = ${TEST_HOST}

.PHONY: all clean $(PROG_NAME) input_in_s3

all: $(PROG_NAME)

$(PROG_NAME):
	go build -o $(BIN_DIR)/$@ .

input_in_s3:
	aws s3 cp ./metamorphosis.txt s3://${BUCKET}

input_in_fastfs:
	aws s3 cp ./metamorphosis.txt s3://${BUCKET}/mapreduce/

input_in_fastfs_medium:
	s3-copy-once ./parking-citations-2M.csv s3://${BUCKET}/mapreduce/

input_in_fastfs_large:
    s3-copy-once ./parking-citations.csv s3://${BUCKET}/mapreduce/

test_wc_local: $(PROG_NAME)
	$(BIN_DIR)/$(PROG_NAME) metamorphosis.txt

test_wc_local_medium: $(PROG_NAME)
	$(BIN_DIR)/$(PROG_NAME) parking-citations-500k.csv

test_wc_s3: $(PROG_NAME) input_in_s3
	$(BIN_DIR)/$(PROG_NAME) --out s3://${BUCKET}/ s3://${BUCKET}/metamorphosis.txt

test_wc_lambda: $(PROG_NAME) input_in_s3
	$(BIN_DIR)/$(PROG_NAME) --lambda --out s3://${BUCKET}/ s3://${BUCKET}/metamorphosis.txt

test_wc_fastfs: $(PROG_NAME) input_in_fastfs
	$(BIN_DIR)/$(PROG_NAME) --out http://localhost:8100/mapreduce http://localhost:8100/mapreduce/metamorphosis.txt

test_wc_s3_medium: $(PROG_NAME) input_in_fastfs_medium
	$(BIN_DIR)/$(PROG_NAME) --out s3://${BUCKET}/mapreduce s3://${BUCKET}/mapreduce/parking-citations-2M.csv

test_wc_fastfs_medium: $(PROG_NAME) input_in_fastfs_medium
	$(BIN_DIR)/$(PROG_NAME) --out http://localhost:8100/mapreduce http://localhost:8100/mapreduce/parking-citations-2M.csv

test_wc_fastfs_lambda: $(PROG_NAME) input_in_fastfs
	$(BIN_DIR)/$(PROG_NAME) --lambda --out http://${HOST}/mapreduce http://${HOST}/mapreduce/metamorphosis.txt

test_wc_lambda_medium: $(PROG_NAME) input_in_fastfs_medium
	$(BIN_DIR)/$(PROG_NAME) --lambda --out s3://${BUCKET}/mapreduce s3://${BUCKET}/mapreduce/parking-citations-2M.csv

test_wc_fastfs_lambda_medium: $(PROG_NAME) input_in_fastfs_medium
	$(BIN_DIR)/$(PROG_NAME) --lambda --out http://${HOST}/mapreduce http://${HOST}/mapreduce/parking-citations-2M.csv

test_wc_lambda_large: $(PROG_NAME) input_in_fastfs_large
	$(BIN_DIR)/$(PROG_NAME) --lambda --out s3://${BUCKET}/mapreduce s3://${BUCKET}/mapreduce/parking-citations.csv

test_wc_fastfs_lambda_large: $(PROG_NAME) input_in_fastfs_large
	$(BIN_DIR)/$(PROG_NAME) --lambda --out http://${HOST}/mapreduce http://${HOST}/mapreduce/parking-citations.csv

# xlarge is already copied on S3 and is 4x size of parking-citations.csv
test_wc_lambda_xlarge: $(PROG_NAME)
	$(BIN_DIR)/$(PROG_NAME) --lambda --out http://${HOST}/mapreduce http://${HOST}/mapreduce/parking-citations4.csv

test_wc_fastfs_lambda_xlarge: $(PROG_NAME)
	$(BIN_DIR)/$(PROG_NAME) --lambda --out http://${HOST}/mapreduce http://${HOST}/mapreduce/parking-citations4.csv


clean:
	find . -name "*.out" -print0 | xargs -0 rm
	rm -f $(BIN_DIR)/$(PROG_NAME) output*
	aws s3 rm s3://${BUCKET} --recursive
