package polar

import (
	"bytes"
	"embed"
	"fmt"
	"testing"
)

//go:embed test_data/*
var testData embed.FS

func TestWriteSimple(t *testing.T) {
	// PolarChunk chunk1, chunk2;
	//        {
	//            var sections = new PolarSection[5];
	//            Arrays.fill(sections, new PolarSection());
	//            chunk1 = new PolarChunk(
	//                    0, 0,
	//                    sections,
	//                    List.of(),
	//                    null,
	//                    new byte[0]
	//            );
	//        }
	//        {
	//            var sections = new PolarSection[5];
	//            Arrays.fill(sections, new PolarSection());
	//            chunk2 = new PolarChunk(
	//                    -1, -1,
	//                    sections,
	//                    List.of(),
	//                    null,
	//                    new byte[0]
	//            );
	//        }
	world := &World{
		Version:     LatestVersion,
		Compression: CompressionNone,
		MinSection:  0,
		MaxSection:  4,
		chunks: map[ChunkIndex]*Chunk{
			ChunkIndexFromXZ(0, 0): {
				X:        0,
				Z:        0,
				Sections: []*Section{},
				BlockEntities: []*BlockEntity{{
					ChunkPos: ChunkPosFromXYZ(3, 4, 5),
					ID:       "minecraft:chest",
					Data: map[string]interface{}{
						"x": int32(3),
					},
				}},
			},
			ChunkIndexFromXZ(-1, -1): {
				X:             -1,
				Z:             -1,
				Sections:      []*Section{},
				BlockEntities: nil,
			},
		},
	}

	data, err := Write(world)
	if err != nil {
		t.Fatal(err)
	}

	javaData, err := testData.ReadFile("test_data/simple_java.polar")
	if err != nil {
		t.Fatal(err)
	}

	// print hex of data and javadata
	fmt.Printf("%x\n", data)
	fmt.Printf("%x\n", javaData)

	if !bytes.Equal(data, javaData) {
		t.Fatal("data mismatch")
	}
}
