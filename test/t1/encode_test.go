package main

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
	"test/t1/gogo"
	"test/t1/google"
)

var data []byte
var rr = &google.R1{
	Uint64: 2065657434543,
	Uint32: 156547,
}

func init() {
	d, err := proto.Marshal(rr)
	if err != nil {
		panic(err)
	}
	data = d
}

func BenchmarkGoogle(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		_ = proto.Unmarshal(data, r)
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}

func BenchmarkGogo(b *testing.B) {
	r := &gogo.R1{}
	for i := 0; i < b.N; i++ {
		_ = r.Unmarshal(data)
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}

func BenchmarkSelf(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		var index int
		l := len(data)

		for j := 0; j < 2; j++ {
			var fieldNum int32
			if l > index {
				fieldNum |= int32(data[index] & 127)
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if l > index {
				fieldNum |= int32(data[index]&127) << 7
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if l > index {
				fieldNum |= int32(data[index]&127) << 14
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if l > index {
				fieldNum |= int32(data[index]&127) << 21
				if data[index] < 128 {
					goto end
				}
			}
			panic(`ErrIntOverflowProto`)
		end:
			index++
			switch fieldNum >> 3 {
			case 4:
				var value uint64
				if l > index {
					value |= uint64(data[index] & 127)
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 7
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 14
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 21
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 28
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 35
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 42
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 49
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 56
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if l > index {
					value |= uint64(data[index]&127) << 63
					if data[index] < 128 {
						goto end2
					}
				}
				panic(`ErrIntOverflowProto`)
			end2:
				index++
				r.Uint64 = value
			case 200:
				var value uint32
				if l > 0 {
					value |= uint32(data[index] & 127)
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if l > 1 {
					value |= uint32(data[index]&127) << 7
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if l > 2 {
					value |= uint32(data[index]&127) << 14
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if l > 3 {
					value |= uint32(data[index]&127) << 21
					if data[index] < 128 {
						goto end3
					}
				}
				panic(`ErrIntOverflowProto`)
			end3:
				index++
				r.Uint32 = value
			}
		}
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}
func BenchmarkSelf2(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		var index int
		for j := 0; j < 2; j++ {
			var fieldNum int32
			if len(data) > index {
				fieldNum |= int32(data[index] & 127)
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 7
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 14
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 21
				if data[index] < 128 {
					goto end
				}
			}
			panic(`ErrIntOverflowProto`)
		end:
			index++
			switch fieldNum >> 3 {
			case 4:
				var value uint64
				if len(data) > index {
					value |= uint64(data[index] & 127)
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 7
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 14
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 21
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 28
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 35
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 42
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 49
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 56
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 63
					if data[index] < 128 {
						goto end2
					}
				}
				panic(`ErrIntOverflowProto`)
			end2:
				index++
				r.Uint64 = value
			case 200:
				var value uint32
				if len(data) > 0 {
					value |= uint32(data[index] & 127)
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 1 {
					value |= uint32(data[index]&127) << 7
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 2 {
					value |= uint32(data[index]&127) << 14
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 3 {
					value |= uint32(data[index]&127) << 21
					if data[index] < 128 {
						goto end3
					}
				}
				panic(`ErrIntOverflowProto`)
			end3:
				index++
				r.Uint32 = value
			}
		}
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}
func BenchmarkSelf3(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		var index int
		for len(data) > index {
			fieldNum := int32(data[index] & 127)
			if data[index] < 128 {
				goto end
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 7
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 14
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 21
				if data[index] < 128 {
					goto end
				}
			}
			panic(`ErrIntOverflowProto`)
		end:
			index++
			switch fieldNum >> 3 {
			case 4:
				var value uint64
				if len(data) > index {
					value |= uint64(data[index] & 127)
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 7
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 14
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 21
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 28
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 35
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 42
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 49
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 56
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 63
					if data[index] < 128 {
						goto end2
					}
				}
				panic(`ErrIntOverflowProto`)
			end2:
				index++
				r.Uint64 = value
			case 200:
				var value uint32
				if len(data) > 0 {
					value |= uint32(data[index] & 127)
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 1 {
					value |= uint32(data[index]&127) << 7
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 2 {
					value |= uint32(data[index]&127) << 14
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 3 {
					value |= uint32(data[index]&127) << 21
					if data[index] < 128 {
						goto end3
					}
				}
				panic(`ErrIntOverflowProto`)
			end3:
				index++
				r.Uint32 = value
			}
		}
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}
func BenchmarkSelf4(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		var index int
		for len(data) > index {
			fieldNum := int32(data[index] & 127)
			if data[index] < 128 {
				goto end
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 7
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 14
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 21
				if data[index] < 128 {
					goto end
				}
			}
			panic(`ErrIntOverflowProto`)
		end:
			index++
			switch fieldNum >> 3 {
			case 4:
				var value uint64
				if len(data) > index {
					value |= uint64(data[index] & 127)
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 7
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 14
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 21
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 28
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 35
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 42
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 49
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 56
					if data[index] < 128 {
						goto end2
					}
				}
				index++
				if len(data) > index {
					value |= uint64(data[index]&127) << 63
					if data[index] < 128 {
						goto end2
					}
				}
				panic(`ErrIntOverflowProto`)
			end2:
				index++
				r.Uint64 = value
			case 200:
				var value uint32
				if len(data) > 0 {
					value |= uint32(data[index] & 127)
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 1 {
					value |= uint32(data[index]&127) << 7
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 2 {
					value |= uint32(data[index]&127) << 14
					if data[index] < 128 {
						goto end3
					}
				}
				index++
				if len(data) > 3 {
					value |= uint32(data[index]&127) << 21
					if data[index] < 128 {
						goto end3
					}
				}
				panic(`ErrIntOverflowProto`)
			end3:
				index++
				r.Uint32 = value
			}
		}
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}

func BenchmarkSelfIf(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		var index int
		l := len(data)

		for j := 0; j < 2; j++ {
			var fieldNum int32
			if l > index {
				fieldNum |= int32(data[index] & 127)
				if data[index] > 127 {
					index++
					if l > index {
						fieldNum |= int32(data[index]&127) << 7
						if data[index] > 127 {
							index++
							if l > index {
								fieldNum |= int32(data[index]&127) << 14
								if data[index] > 127 {
									index++
									if l > index {
										fieldNum |= int32(data[index]&127) << 21
										if data[index] > 127 {
											panic(`ErrIntOverflowProto`)
										}
									}
								}
							}
						}
					}
				}
			}
			index++
			switch fieldNum >> 3 {
			case 4:
				var value uint64
				if l > index {
					value |= uint64(data[index] & 127)
					if data[index] > 127 {
						index++
						if l > index {
							value |= uint64(data[index]&127) << 7
							if data[index] > 127 {
								index++
								if l > index {
									value |= uint64(data[index]&127) << 14
									if data[index] > 127 {
										index++
										if l > index {
											value |= uint64(data[index]&127) << 21
											if data[index] > 127 {
												index++
												if l > index {
													value |= uint64(data[index]&127) << 28
													if data[index] > 127 {
														index++
														if l > index {
															value |= uint64(data[index]&127) << 35
															if data[index] > 127 {
																index++
																if l > index {
																	value |= uint64(data[index]&127) << 42
																	if data[index] > 127 {
																		index++
																		if l > index {
																			value |= uint64(data[index]&127) << 49
																			if data[index] > 127 {
																				index++
																				if l > index {
																					value |= uint64(data[index]&127) << 56
																					if data[index] > 127 {
																						index++
																						if l > index {
																							value |= uint64(data[index]&127) << 63
																							if data[index] > 127 {
																								panic(`ErrIntOverflowProto`)
																							}
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
				index++
				r.Uint64 = value
			case 200:
				var value uint32
				if l > index {
					value |= uint32(data[index] & 127)
					if data[index] > 127 {
						index++
						if l > index {
							value |= uint32(data[index]&127) << 7
							if data[index] > 127 {
								index++
								if l > index {
									value |= uint32(data[index]&127) << 14
									if data[index] > 127 {
										index++
										if l > 3 {
											value |= uint32(data[index]&127) << 21
											if data[index] < 128 {
												panic(`ErrIntOverflowProto`)
											}
										}
									}
								}
							}
						}
					}
				}
				index++
				r.Uint32 = value
			}
		}
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}
func BenchmarkSelfLen(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		var index int
		for len(data) > index {
			fieldNum := int32(data[index] & 127)
			if data[index] > 127 {
				index++
				if len(data) > index {
					fieldNum |= int32(data[index]&127) << 7
					if data[index] > 127 {
						index++
						if len(data) > index {
							fieldNum |= int32(data[index]&127) << 14
							if data[index] > 127 {
								index++
								if len(data) > index {
									fieldNum |= int32(data[index]&127) << 21
									if data[index] > 127 {
										panic(`ErrIntOverflowProto`)
									}
								}
							}
						}
					}
				}
			}
			index++
			switch fieldNum >> 3 {
			case 4:
				var value uint64
				if len(data) > index {
					value |= uint64(data[index] & 127)
					if data[index] > 127 {
						index++
						if len(data) > index {
							value |= uint64(data[index]&127) << 7
							if data[index] > 127 {
								index++
								if len(data) > index {
									value |= uint64(data[index]&127) << 14
									if data[index] > 127 {
										index++
										if len(data) > index {
											value |= uint64(data[index]&127) << 21
											if data[index] > 127 {
												index++
												if len(data) > index {
													value |= uint64(data[index]&127) << 28
													if data[index] > 127 {
														index++
														if len(data) > index {
															value |= uint64(data[index]&127) << 35
															if data[index] > 127 {
																index++
																if len(data) > index {
																	value |= uint64(data[index]&127) << 42
																	if data[index] > 127 {
																		index++
																		if len(data) > index {
																			value |= uint64(data[index]&127) << 49
																			if data[index] > 127 {
																				index++
																				if len(data) > index {
																					value |= uint64(data[index]&127) << 56
																					if data[index] > 127 {
																						index++
																						if len(data) > index {
																							value |= uint64(data[index]&127) << 63
																							if data[index] > 127 {
																								panic(`ErrIntOverflowProto`)
																							}
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
				index++
				r.Uint64 = value
			case 200:
				var value uint32
				if len(data) > index {
					value |= uint32(data[index] & 127)
					if data[index] > 127 {
						index++
						if len(data) > index {
							value |= uint32(data[index]&127) << 7
							if data[index] > 127 {
								index++
								if len(data) > index {
									value |= uint32(data[index]&127) << 14
									if data[index] > 127 {
										index++
										if len(data) > index {
											value |= uint32(data[index]&127) << 21
											if data[index] < 128 {
												panic(`ErrIntOverflowProto`)
											}
										}
									}
								}
							}
						}
					}
				}
				index++
				r.Uint32 = value
			}
		}
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}
func BenchmarkSelfLen2(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		var index int
		for len(data) > index {
			var fieldNum int32
			if data[index] > 127 {
				if len(data) > index+1 {
					if data[index+1] > 127 {
						if len(data) > index+2 {
							if data[index+2] > 127 {
								if len(data) > index+3 {
									if data[index+3] > 127 {
										panic(`ErrIntOverflowProto`)
									} else {
										fieldNum = int32(data[index]&127) |
											int32(data[index+1]&127)<<7 |
											int32(data[index+2]&127)<<14 |
											int32(data[index+3]&127)<<21
										index += 4
									}
								} else {
									panic(`---`)
								}
							} else {
								fieldNum = int32(data[index]&127) |
									int32(data[index+1]&127)<<7 |
									int32(data[index+2]&127)<<14
								index += 3
							}
						} else {
							panic(`---`)
						}
					} else {
						fieldNum = int32(data[index]&127) |
							int32(data[index+1]&127)<<7
						index += 2
					}
				} else {
					panic(`---`)
				}
			} else {
				fieldNum = int32(data[index] & 127)
				index += 1
			}
			switch fieldNum >> 3 {
			case 4:
				var value uint64
				if len(data) > index {
					if data[index] > 127 {
						if len(data) > index+1 {
							if data[index+1] > 127 {
								if len(data) > index+2 {
									if data[index+2] > 127 {
										if len(data) > index+3 {
											if data[index+3] > 127 {
												if len(data) > index+4 {
													if data[index+4] > 127 {
														if len(data) > index+5 {
															if data[index+5] > 127 {
																if len(data) > index+6 {
																	if data[index+6] > 127 {
																		if len(data) > index+7 {
																			if data[index+7] > 127 {
																				if len(data) > index+8 {
																					if data[index+8] > 127 {
																						if len(data) > index+9 {
																							if data[index+9] > 127 {
																								panic(`ErrIntOverflowProto`)
																							} else {
																								value = uint64(data[index]&127) |
																									uint64(data[index+1]&127)<<7 |
																									uint64(data[index+2]&127)<<14 |
																									uint64(data[index+3]&127)<<21 |
																									uint64(data[index+4]&127)<<28 |
																									uint64(data[index+5]&127)<<35 |
																									uint64(data[index+6]&127)<<42 |
																									uint64(data[index+7]&127)<<49 |
																									uint64(data[index+8]&127)<<56 |
																									uint64(data[index+9]&127)<<63
																								index += 10
																							}
																						} else {
																							panic(`---`)
																						}
																					} else {
																						value = uint64(data[index]&127) |
																							uint64(data[index+1]&127)<<7 |
																							uint64(data[index+2]&127)<<14 |
																							uint64(data[index+3]&127)<<21 |
																							uint64(data[index+4]&127)<<28 |
																							uint64(data[index+5]&127)<<35 |
																							uint64(data[index+6]&127)<<42 |
																							uint64(data[index+7]&127)<<49 |
																							uint64(data[index+8]&127)<<56
																						index += 9
																					}
																				} else {
																					panic(`---`)
																				}
																			} else {
																				value = uint64(data[index]&127) |
																					uint64(data[index+1]&127)<<7 |
																					uint64(data[index+2]&127)<<14 |
																					uint64(data[index+3]&127)<<21 |
																					uint64(data[index+4]&127)<<28 |
																					uint64(data[index+5]&127)<<35 |
																					uint64(data[index+6]&127)<<42 |
																					uint64(data[index+7]&127)<<49
																				index += 8
																			}
																		} else {
																			panic(`---`)
																		}
																	} else {
																		value = uint64(data[index]&127) |
																			uint64(data[index+1]&127)<<7 |
																			uint64(data[index+2]&127)<<14 |
																			uint64(data[index+3]&127)<<21 |
																			uint64(data[index+4]&127)<<28 |
																			uint64(data[index+5]&127)<<35 |
																			uint64(data[index+6]&127)<<42
																		index += 7
																	}
																} else {
																	panic(`---`)
																}
															} else {
																value = uint64(data[index]&127) |
																	uint64(data[index+1]&127)<<7 |
																	uint64(data[index+2]&127)<<14 |
																	uint64(data[index+3]&127)<<21 |
																	uint64(data[index+4]&127)<<28 |
																	uint64(data[index+5]&127)<<35
																index += 6
															}
														} else {
															panic(`---`)
														}
													} else {
														value = uint64(data[index]&127) |
															uint64(data[index+1]&127)<<7 |
															uint64(data[index+2]&127)<<14 |
															uint64(data[index+3]&127)<<21 |
															uint64(data[index+4]&127)<<28
														index += 5
													}
												} else {
													panic(`---`)
												}
											} else {
												value = uint64(data[index]&127) |
													uint64(data[index+1]&127)<<7 |
													uint64(data[index+2]&127)<<14 |
													uint64(data[index+3]&127)<<21
												index += 4
											}
										} else {
											panic(`---`)
										}
									} else {
										value = uint64(data[index]&127) |
											uint64(data[index+1]&127)<<7 |
											uint64(data[index+2]&127)<<14
										index += 3
									}
								} else {
									panic(`---`)
								}
							} else {
								value = uint64(data[index]&127) |
									uint64(data[index+1]&127)<<7
								index += 2
							}
						} else {
							panic(`---`)
						}
					} else {
						value = uint64(data[index] & 127)
						index += 1
					}
				} else {
					panic(`---`)
				}
				r.Uint64 = value
			case 200:
				var value uint32
				if len(data) > index {
					value |= uint32(data[index] & 127)
					if data[index] > 127 {
						index++
						if len(data) > index {
							value |= uint32(data[index]&127) << 7
							if data[index] > 127 {
								index++
								if len(data) > index {
									value |= uint32(data[index]&127) << 14
									if data[index] > 127 {
										index++
										if len(data) > index {
											value |= uint32(data[index]&127) << 21
											if data[index] < 128 {
												panic(`ErrIntOverflowProto`)
											}
										}
									}
								}
							}
						}
					}
				}
				index++
				r.Uint32 = value
			}
		}
		if rr.Uint32 != r.Uint32 || rr.Uint64 != r.Uint64 {
			panic(`eq`)
		}
	}
}

func BenchmarkLen(b *testing.B) {
	qrr := []int{1, 2, 3, 4, 5, 6, 7}
	for i := 0; i < b.N; i++ {
		_ = len(qrr) > i
	}
}
func BenchmarkIf(b *testing.B) {
	qrr := []int{1, 2, 3, 4, 5, 6, 7}
	l := len(qrr)
	for i := 0; i < b.N; i++ {
		_ = l > i
	}
}

func Test(*testing.T) {
	r := &google.R1{
		Uint64: 0,
		Uint32: 0,
		Int64:  1,
		Int32:  1,
		Sint64: 0,
		Sint32: 0,
	}
	data, err := proto.Marshal(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%b\n", data)
	fmt.Printf("%d\n", data)
}
