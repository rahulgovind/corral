BUCKET = ${AWS_TEST_BUCKET}
BIN_DIR = ./bin
PROG_NAME = word_count

.PHONY: all clean $(PROG_NAME) input_in_s3

all: $(PROG_NAME)

$(PROG_NAME):
	go build -o $(BIN_DIR)/$@ .

input_in_s3:
	aws s3 cp ./metamorphosis.txt s3://${BUCKET}

input_in_fastfs:
	aws s3 cp ./metamorphosis.txt s3://${BUCKET}/mapreduce/

input_in_fastfs_medium:
	s3-copy-once ./parking-citations-500k.csv s3://${BUCKET}/mapreduce/

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

test_wc_fastfs_medium: $(PROG_NAME) input_in_fastfs_medium
	$(BIN_DIR)/$(PROG_NAME) --out http://localhost:8100/mapreduce http://localhost:8100/mapreduce/parking-citations-500k.csv

clean:
	find . -name "*.out" -print0 | xargs -0 rm
	rm -f $(BIN_DIR)/$(PROG_NAME) output*
	aws s3 rm s3://${BUCKET} --recursive
