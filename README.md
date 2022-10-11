<h1 align="center">Brain 🧠</h1>
<p>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> Your external brain from the terminal.

## What is this?

Brain is a command line tool for knowledge workers that aids with information storage and retrieval. Or, in other words, a fancy note-taking app.

## Sounds mint, how do I use it?

So you're reading a technical document, or in a meeting with your tech lead that is getting a little complicated, and you want to offload some of that information so you can nod like you understand what's being said.

From your terminal, start a session by typing:

```sh
$ brain
```

This will spawn an editor (configurable by $EDITOR) which you can start writing in. Think of this more as a raw brain dump than a structured document. Try to use as many keywords that will help with retrieval later on. The goal is to optimize for fast retrieval instead of the correctness of the note, your own brain will correct that for you when you re-read the note later on.

As an example, let's say we're reading about Apache Kafka... instead of writing

"Kafka is a distributed event store and stream-processing platform."

We could instead write

"Kafka is a distributed [scalable] event [message] store and stream [message, event] processing platform [platform]".

What we're doing here is writing our notes normally, but adding ..

Remember, the goal is to optimize for fast retrieval. What this means in practice is that when we want to recall what we
wrote about Kafka, we don't have to memorize a single concept. The cost of writing more is amortized because we
will write once and read many times. The more nuanced the concept the more you will benefit from this approach.

Ok, now comes the fun part. You wrote some notes last week when you were learning about distributed systems, and you want
to remember what your wrote about the CAP theorem.

```shell
brain r CAP
```

```shell
...
...
...
```

We've been reading quite a lot lately it seems.. we're going to have to drill down on what we really want


## Install

```sh
go get github.com/sno6/brain
```

## Usage

```sh
Usage:
  brain [flags]

Flags:
  -h, --help                 help for reaper
      --mh int               Only allow images with a height >= to this value (default -1)
      --mw int               Only allow images with a width >= to this value (default -1)
  -o, --out string           Output directory (default "./")
  -s, --search stringArray   Search terms
  -t, --timeout duration     HTTP request timeout in seconds (default 15s)
  -v, --verbose              Log errors
```


## Details

Configuration:

- Editor 'vim'
- Keyword char '['

File structure:

Query patterns:

Users can search by datetime range
Users can search by keywords 
Users can search by direct string match
Users can search by fuzzy string match

Writing patterns:

Users write to a flat file where the latest note is written to the end of the file. This op should be constant or near.
Users can edit existing notes although this is discouraged and doesn't have to be optimized.
Users can delete existing notes although this is discouraged and doesn't have to be optimized.



