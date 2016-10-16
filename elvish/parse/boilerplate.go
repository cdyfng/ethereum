package parse

func (n *Chunk) addToPipelines(ch *Pipeline) {
	n.Pipelines = append(n.Pipelines, ch)
	addChild(n, ch)
}

func parseChunk(ps *parser) *Chunk {
	n := &Chunk{node: node{begin: ps.pos}}
	n.parse(ps)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Pipeline) addToForms(ch *Form) {
	n.Forms = append(n.Forms, ch)
	addChild(n, ch)
}

func parsePipeline(ps *parser) *Pipeline {
	n := &Pipeline{node: node{begin: ps.pos}}
	n.parse(ps)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Form) addToAssignments(ch *Assignment) {
	n.Assignments = append(n.Assignments, ch)
	addChild(n, ch)
}

func (n *Form) setControl(ch *Control) {
	n.Control = ch
	addChild(n, ch)
}

func (n *Form) setHead(ch *Compound) {
	n.Head = ch
	addChild(n, ch)
}

func (n *Form) addToArgs(ch *Compound) {
	n.Args = append(n.Args, ch)
	addChild(n, ch)
}

func (n *Form) addToOpts(ch *MapPair) {
	n.Opts = append(n.Opts, ch)
	addChild(n, ch)
}

func (n *Form) addToRedirs(ch *Redir) {
	n.Redirs = append(n.Redirs, ch)
	addChild(n, ch)
}

func (n *Form) setExitusRedir(ch *ExitusRedir) {
	n.ExitusRedir = ch
	addChild(n, ch)
}

func parseForm(ps *parser) *Form {
	n := &Form{node: node{begin: ps.pos}}
	n.parse(ps)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Assignment) setDst(ch *Indexing) {
	n.Dst = ch
	addChild(n, ch)
}

func (n *Assignment) setSrc(ch *Compound) {
	n.Src = ch
	addChild(n, ch)
}

func parseAssignment(ps *parser) *Assignment {
	n := &Assignment{node: node{begin: ps.pos}}
	n.parse(ps)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Control) setCondition(ch *Chunk) {
	n.Condition = ch
	addChild(n, ch)
}

func (n *Control) setIterator(ch *Indexing) {
	n.Iterator = ch
	addChild(n, ch)
}

func (n *Control) setArray(ch *Array) {
	n.Array = ch
	addChild(n, ch)
}

func (n *Control) setBody(ch *Chunk) {
	n.Body = ch
	addChild(n, ch)
}

func (n *Control) addToConditions(ch *Chunk) {
	n.Conditions = append(n.Conditions, ch)
	addChild(n, ch)
}

func (n *Control) addToBodies(ch *Chunk) {
	n.Bodies = append(n.Bodies, ch)
	addChild(n, ch)
}

func (n *Control) setElseBody(ch *Chunk) {
	n.ElseBody = ch
	addChild(n, ch)
}

func (n *Control) setExceptBody(ch *Chunk) {
	n.ExceptBody = ch
	addChild(n, ch)
}

func (n *Control) setExceptVar(ch *Indexing) {
	n.ExceptVar = ch
	addChild(n, ch)
}

func (n *Control) setFinallyBody(ch *Chunk) {
	n.FinallyBody = ch
	addChild(n, ch)
}

func parseControl(ps *parser, leader string) *Control {
	n := &Control{node: node{begin: ps.pos}}
	n.parse(ps, leader)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *ExitusRedir) setDest(ch *Compound) {
	n.Dest = ch
	addChild(n, ch)
}

func parseExitusRedir(ps *parser) *ExitusRedir {
	n := &ExitusRedir{node: node{begin: ps.pos}}
	n.parse(ps)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Redir) setDest(ch *Compound) {
	n.Dest = ch
	addChild(n, ch)
}

func (n *Redir) setSource(ch *Compound) {
	n.Source = ch
	addChild(n, ch)
}

func parseRedir(ps *parser, dest *Compound) *Redir {
	n := &Redir{node: node{begin: ps.pos}}
	n.parse(ps, dest)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Compound) addToIndexings(ch *Indexing) {
	n.Indexings = append(n.Indexings, ch)
	addChild(n, ch)
}

func parseCompound(ps *parser, head bool) *Compound {
	n := &Compound{node: node{begin: ps.pos}}
	n.parse(ps, head)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Indexing) setHead(ch *Primary) {
	n.Head = ch
	addChild(n, ch)
}

func (n *Indexing) addToIndicies(ch *Array) {
	n.Indicies = append(n.Indicies, ch)
	addChild(n, ch)
}

func parseIndexing(ps *parser, head bool) *Indexing {
	n := &Indexing{node: node{begin: ps.pos}}
	n.parse(ps, head)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Array) addToCompounds(ch *Compound) {
	n.Compounds = append(n.Compounds, ch)
	addChild(n, ch)
}

func parseArray(ps *parser, allowSemicolon bool) *Array {
	n := &Array{node: node{begin: ps.pos}}
	n.parse(ps, allowSemicolon)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *Primary) setList(ch *Array) {
	n.List = ch
	addChild(n, ch)
}

func (n *Primary) setChunk(ch *Chunk) {
	n.Chunk = ch
	addChild(n, ch)
}

func (n *Primary) addToMapPairs(ch *MapPair) {
	n.MapPairs = append(n.MapPairs, ch)
	addChild(n, ch)
}

func (n *Primary) addToBraced(ch *Compound) {
	n.Braced = append(n.Braced, ch)
	addChild(n, ch)
}

func parsePrimary(ps *parser, head bool) *Primary {
	n := &Primary{node: node{begin: ps.pos}}
	n.parse(ps, head)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}

func (n *MapPair) setKey(ch *Compound) {
	n.Key = ch
	addChild(n, ch)
}

func (n *MapPair) setValue(ch *Compound) {
	n.Value = ch
	addChild(n, ch)
}

func parseMapPair(ps *parser) *MapPair {
	n := &MapPair{node: node{begin: ps.pos}}
	n.parse(ps)
	n.end = ps.pos
	n.sourceText = ps.src[n.begin:n.end]
	return n
}
