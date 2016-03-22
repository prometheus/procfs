package procfs

type Diskstats struct {
}

func NewMeminfo() (Meminfo, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Meminfo{}, err
	}

	return fs.NewMeminfo()
}

// NewMeminfo returns an information about current kernel/system statistics.
func (fs FS) NewMeminfo() (m Meminfo, err error) {
	f, err := os.Open(fs.Path("meminfo"))
	if err != nil {
		return Meminfo{}, err
	}
	defer f.Close()

	st := reflect.TypeOf(m)
	re := regexp.MustCompile(m.regex())
	s := bufio.NewScanner(f)

	for s.Scan() {
		line := s.Text()

		submatch := re.FindAllStringSubmatch(line, 1)
		if submatch == nil {
			continue
		}

		key := submatch[0][1]
		val := submatch[0][2]

		for i := 0; i < st.NumField(); i++ {
			field := st.Field(i)
			if field.Tag.Get("meminfo") == key {
				v, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					// no op
				}
				reflect.ValueOf(&m).Elem().Field(i).SetInt(v)
			}
		}
	}

	return m, nil
}

func (m Meminfo) regex() string {
	return "([A-Za-z0-9()_]*): *([0-9]*).*$"
}

/*
# collect i/o load
if (open(FILE, "/proc/diskstats")) {
    while (my $line = <FILE>) {
        $line =~ s/^\s+|\s+$//g;
        my @cols = split(/\s+/, $line);
        splice(@cols, 0, 3);
        next if (scalar(@cols) != 11);
        $load{d_io_reads} += $cols[0];
        $load{d_io_read_sectors} += $cols[2];
        $load{d_io_read_time} += $cols[3];
        $load{d_io_writes} += $cols[4];
        $load{d_io_write_sectors} += $cols[6];
        $load{d_io_write_time} += $cols[7];
    }
    close FILE;
    $load{d_io_ops} = $load{d_io_reads} + $load{d_io_writes};
    $load{d_io_sectors} = $load{d_io_read_sectors} + $load{d_io_write_sectors};
    $load{d_io_time} = $load{d_io_read_time} + $load{d_io_write_time};
}

*/
